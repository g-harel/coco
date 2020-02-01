package npm

import (
	"fmt"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/flags"
	"github.com/g-harel/coco/internal/table"
)

var _ collectors.Collector = &Collector{}

// Collector satisfies the collector interface to fetch and
// display NPM package download info.
type Collector struct {
	packages []*pkg
}

// Collect fetches packages from all owners in parallel.
func (c *Collector) Collect(h collectors.ErrorHandler) {
	exec.ParallelN(len(flags.NpmOwners), func(n int) {
		handleOwner(func(r *pkg, err error) {
			if err != nil {
				h(err)
				return
			}
			exec.Safe(func() {
				c.packages = append(c.packages, r)
			})
		}, flags.NpmOwners[n])
	})
}

// Format creates a table from the collected download data.
// It allows the shown packages to be filtered by weekly
// downloads.
func (c *Collector) Format() string {
	if len(c.packages) == 0 {
		return ""
	}

	t := table.Table{}

	owners := ""
	for i := 0; i < len(flags.NpmOwners); i++ {
		owners += " " + flags.NpmOwners[i]
	}
	t.Title(fmt.Sprintf("Npm package downloads |%v", owners))

	t.Headers(
		"PACKAGE",
		"DOWNLOADS",
		"TOTAL",
		"LINK",
	)

	for i := 0; i < len(c.packages); i++ {
		p := c.packages[i]
		if p.Weekly < *flags.NpmWeekly {
			continue
		}
		link := "https://npmjs.com/package/" + p.Name
		t.Add(
			p.Name,
			p.Weekly,
			p.Total,
			link,
		)

	}
	t.Sort(1, 2)

	return t.String()
}
