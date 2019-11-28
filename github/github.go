package github

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/g-harel/coco/internal"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is used to configure communication with the GitHub API.
type Client struct {
	client *github.Client
}

// repositoryStats represents metadata about a single repository.
type repositoryStats struct {
	Name   string
	Owner  string
	Views  int
	Today  int
	Unique int
	Error  error
}

// repositoryStatsList represents a slice of repositories.
type repositoryStatsList []*repositoryStats

// Filter returns a slice of the repositories that pass the test function.
func (repos repositoryStatsList) Filter(test func(*repositoryStats) bool) repositoryStatsList {
	kept := repositoryStatsList{}
	for _, r := range repos {
		if test(r) {
			kept = append(kept, r)
		}
	}
	return kept
}

// String returns a formatted view of the repository data.
func (repos repositoryStatsList) String() string {
	t := internal.Table{}
	t.Headers("REPO", "VIEWS", "DAY", "UNIQUE", "LINK")
	for _, r := range repos {
		t.Add(
			r.Name,
			r.Views,
			r.Today,
			r.Unique,
			"https://github.com/"+r.Owner+"/"+r.Name+"/graphs/traffic",
		)
	}
	return t.Format(1, 2, 3)
}

// NewClient creates and configures a new Client.
func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &Client{
		github.NewClient(internal.NewLoggingClient(tc)),
	}
}

// Repositories queries for all repositories linked to the given owners.
func (c *Client) Repositories(users []string) (repositoryStatsList, error) {
	repos := repositoryStatsList{}

	for _, name := range users {
		name = strings.TrimSpace(name)
		list := []*github.Repository{}

		page := 1
		for {
			// Paginated list of repositories is queried.
			r, res, err := c.client.Repositories.List(
				context.Background(),
				strings.TrimSpace(name),
				&github.RepositoryListOptions{
					Type:        "all",
					ListOptions: github.ListOptions{Page: page},
				},
			)
			if err != nil {
				return nil, fmt.Errorf("fetch repositories: %v", err)
			}

			list = append(list, r...)

			// NextPage value of zero indicates there are no additional pages.
			page = res.NextPage
			if page == 0 {
				break
			}
		}

		// Repository data is translated into a
		for _, l := range list {
			repos = append(repos, &repositoryStats{
				Owner: l.GetOwner().GetLogin(),
				Name:  l.GetName(),
			})
		}
	}

	return repos, nil
}

// Traffic queries for repository traffic data in parallel.
// Items in the input slice are modified in place.
func (c *Client) Traffic(repos repositoryStatsList) repositoryStatsList {
	rch := make(chan *repositoryStats, len(repos))

	for i := 0; i < len(repos); i++ {
		go func(i int) {
			traffic, _, err := c.client.Repositories.ListTrafficViews(
				context.Background(),
				repos[i].Owner,
				repos[i].Name,
				nil,
			)
			if err != nil {
				rch <- &repositoryStats{Error: err}
				return
			}

			// Traffic from today is highlighted.
			var trafficToday int
			for _, stat := range traffic.Views {
				if isToday(stat.GetTimestamp()) {
					trafficToday += stat.GetCount()
				}
			}

			rch <- &repositoryStats{
				Name:   repos[i].Name,
				Owner:  repos[i].Owner,
				Views:  traffic.GetCount(),
				Today:  trafficToday,
				Unique: traffic.GetUniques(),
			}
		}(i)
	}

	r := make(repositoryStatsList, len(repos))
	for i := 0; i < len(repos); i++ {
		r[i] = <-rch
	}

	return r
}

func isToday(t github.Timestamp) bool {
	f := "2006-01-02"
	return time.Now().Format(f) == t.Format(f)
}

func Repositories(token string, users ...string) string {
	gh := NewClient(token)

	repos, err := gh.Repositories(users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch repositories: %v", err)
		os.Exit(1)
	}

	// Remove duplicate repositories (usernames might have overlap).
	visited := make(map[string]bool)
	repos = repos.Filter(func(r *repositoryStats) bool {
		if visited[r.Owner+r.Name] {
			return false
		}
		visited[r.Owner+r.Name] = true
		return true
	})

	// Fetch traffic data for all repositories.
	repos = gh.Traffic(repos)

	// Remove repos with errors or no reported views (in the past two weeks).
	repos = repos.Filter(func(r *repositoryStats) bool {
		if r.Error != nil {
			// Fetching errors are swallowed to avoid crowding the output (subject to change).
			return false
		}
		return r.Views != 0
	})

	// Remove repos with very little views.
	repos = repos.Filter(func(r *repositoryStats) bool {
		return !(r.Today < 2 && r.Views < 4)
	})

	return repos.String()
}
