package npm

import (
	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/flags"
	"github.com/g-harel/coco/internal/table"
)

var _ collectors.Collector = &Collector{}

type Collector struct {
	packages []*pkg
}

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

func (c *Collector) Format() string {
	t := table.Table{}
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
