package github

import (
	"fmt"
	"strconv"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/flags"
	"github.com/g-harel/coco/internal/state"
	"github.com/g-harel/coco/internal/table"
)

var _ collectors.Collector = &Collector{}

// Collector satisfies the collector interface to fetch and
// display GitHub repo views info.
type Collector struct {
	repos []*repo
}

// Collect fetches repo data from all owners in parallel.
func (c *Collector) Collect(h collectors.ErrorHandler) {
	if len(flags.GithubOwners) > 0 && *flags.GithubToken == "" {
		h(fmt.Errorf("missing github token"))
		return
	}
	exec.ParallelN(len(flags.GithubOwners), func(n int) {
		handleOwner(func(r *repo, err error) {
			if err != nil {
				h(err)
				return
			}
			exec.Safe(func() {
				c.repos = append(c.repos, r)
			})
		}, *flags.GithubToken, flags.GithubOwners[n])
	})
}

// Format creates a table from the collected views data. It
// allows the shown repos to be filtered by daily views,
// total views and stars.
func (c *Collector) Format() string {
	if len(c.repos) == 0 {
		return ""
	}
	t := table.Table{}
	owners := ""
	for i := 0; i < len(flags.GithubOwners); i++ {
		owners += " " + flags.GithubOwners[i]
	}
	t.Title(fmt.Sprintf("GitHub repo stats |%v", owners))
	t.Headers(
		"REPO",
		"VIEWS",
		"UNIQUE",
		"TODAY",
		"LINK",
	)
	storedData := state.NewFromFile(*flags.StateFile)
	for i := 0; i < len(c.repos); i++ {
		r := c.repos[i]
		url := fmt.Sprintf("https://github.com/%v/%v/graphs/traffic", r.Owner, r.Name)
		newStars := r.Stars - storedData.ReadIntOr(url, 0)
		storedData.Write(url, strconv.Itoa(r.Stars))
		if r.Today < *flags.GithubToday &&
			r.Views < *flags.GithubViews &&
			r.Stars < *flags.GithubStars &&
			newStars < *flags.GithubNewStars {
			continue
		}
		name := r.Name
		if newStars > 0 {
			name += fmt.Sprintf("+%v", newStars)
		}
		t.Add(name, r.Views, r.Unique, r.Today, url)
	}
	storedData.Save()
	t.Sort(1, 3, 2)
	return t.String()
}
