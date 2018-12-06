package main

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Repositories represents a slice of repositories.
type Repositories []*Repository

// Filter returns a slice of the repositories that pass the test function.
func (repos Repositories) Filter(test func(*Repository) bool) Repositories {
	kept := Repositories{}
	for _, r := range repos {
		if test(r) {
			kept = append(kept, r)
		}
	}
	return kept
}

// String returns a formatted view of the repository data.
func (repos Repositories) String() string {
	sort.Sort(repos)

	data := [][]string{}
	for _, r := range repos {
		data = append(data, []string{
			r.Name,
			strconv.Itoa(r.Views),
			strconv.Itoa(r.Today),
			strconv.Itoa(r.Unique),
			"https://github.com/" + r.Owner + "/" + r.Name + "/graphs/traffic",
		})
	}

	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"repo", "views", "day", "unique", "link"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})
	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

// Remaining functions implement sort.Interface.

func (repos Repositories) Len() int {
	return len(repos)
}

func (repos Repositories) Less(i, j int) bool {
	a := repos[i]
	b := repos[j]
	// Sort by: total views -> today's views -> unique views -> name
	if a.Views == b.Views {
		if a.Today == b.Today {
			if a.Unique == b.Unique {
				return strings.Compare(a.Name, b.Name) < 0
			}
			return a.Unique > b.Unique
		}
		return a.Today > b.Today
	}
	return a.Views > b.Views
}

func (repos Repositories) Swap(i, j int) {
	repos[i], repos[j] = repos[j], repos[i]
}
