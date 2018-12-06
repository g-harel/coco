package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is used to configure communication with the GitHub API.
type Client struct {
	client *github.Client
}

// NewClient creates and configures a new Client.
func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &Client{
		github.NewClient(tc),
	}
}

// Repositories queries for all repositories linked to the given owners.
func (c *Client) Repositories(users []string) (Repositories, error) {
	repos := Repositories{}

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
			repos = append(repos, &Repository{
				Owner: l.GetOwner().GetLogin(),
				Name:  l.GetName(),
			})
		}
	}

	return repos, nil
}

// Traffic queries for repository traffic data in parallel.
// Items in the input slice are modified in place.
func (c *Client) Traffic(repos Repositories) Repositories {
	rch := make(chan *Repository, len(repos))

	for i := 0; i < len(repos); i++ {
		go func(i int) {
			traffic, _, err := c.client.Repositories.ListTrafficViews(
				context.Background(),
				repos[i].Owner,
				repos[i].Name,
				nil,
			)
			if err != nil {
				rch <- &Repository{Error: err}
				return
			}

			// Traffic from today is highlighted.
			var trafficToday *github.TrafficData
			for _, stat := range traffic.Views {
				if isToday(stat.GetTimestamp()) {
					trafficToday = stat
					break
				}
			}

			rch <- &Repository{
				Name:   repos[i].Name,
				Owner:  repos[i].Owner,
				Views:  traffic.GetCount(),
				Today:  trafficToday.GetCount(),
				Unique: traffic.GetUniques(),
			}
		}(i)
	}

	r := make(Repositories, len(repos))
	for i := 0; i < len(repos); i++ {
		r[i] = <-rch
	}

	return r
}

func isToday(t github.Timestamp) bool {
	y, m, d := time.Now().Date()
	cy, cm, cd := t.Date()
	return y == cy && m == cm && d == cd
}
