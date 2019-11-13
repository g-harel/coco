package npm

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Pkg struct {
	Name   string
	Weekly int
	Total  int
}

type PackageList []*Pkg

func (packages PackageList) String() string {
	sort.Sort(packages)

	data := [][]string{}
	for _, p := range packages {
		data = append(data, []string{
			p.Name,
			strconv.Itoa(p.Weekly),
			strconv.Itoa(p.Total),
			"https://npmjs.com/package/" + p.Name,
		})
	}

	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"package", "downloads", "total", "link"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})
	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

// Remaining functions implement sort.Interface.

func (packages PackageList) Len() int {
	return len(packages)
}

func (packages PackageList) Less(i, j int) bool {
	a := packages[i]
	b := packages[j]
	// Sort by: weekly downloads -> total downloads -> name
	if a.Weekly == b.Weekly {
		if a.Total == b.Total {
			return strings.Compare(a.Name, b.Name) < 0
		}
		return a.Total > b.Total
	}
	return a.Weekly > b.Weekly
}

func (packages PackageList) Swap(i, j int) {
	packages[i], packages[j] = packages[j], packages[i]
}
