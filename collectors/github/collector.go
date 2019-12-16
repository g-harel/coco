package github

import (
	"fmt"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/flags"
	"github.com/g-harel/coco/internal/table"
)

var _ collectors.Collector = &Collector{}

type Collector struct {
	repos []*repo
}

func (c *Collector) Collect(h collectors.ErrorHandler) {
	exec.ParallelN(len(flags.GithubOwners), func(n int) {
		handleOwner(func(r *repo, err error) {
			if err != nil {
				h(err)
			} else {
				c.repos = append(c.repos, r)
			}
		}, *flags.GithubToken, flags.GithubOwners[n])
	})
}

func (c *Collector) Format() string {
	t := table.Table{}
	t.Headers(
		"REPO",
		"VIEWS",
		"UNIQUE",
		"TODAY",
		"LINK",
	)
	for i := 0; i < len(c.repos); i++ {
		r := c.repos[i]
		if r.Today < *flags.GithubToday &&
			r.Views < *flags.GithubViews &&
			r.Stars < *flags.GithubStars {
			continue
		}
		t.Add(
			fmt.Sprintf("%v*%v", r.Name, r.Stars),
			r.Views,
			r.Unique,
			r.Today,
			fmt.Sprintf("https://github.com/%v/%v/graphs/traffic", r.Owner, r.Name),
		)
	}
	t.Sort(1, 3, 2)
	return t.String()
}
