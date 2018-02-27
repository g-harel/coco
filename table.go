package main

import (
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type TrafficTable struct {
	repos []*Repo
}

func NewTrafficTable(repos []*Repo) *TrafficTable {
	return &TrafficTable{repos}
}

func (t *TrafficTable) Print() {
	sort.Sort(t)

	data := [][]string{}
	for _, r := range t.repos {
		data = append(data, []string{
			r.Name,
			strconv.Itoa(r.Views),
			strconv.Itoa(r.Today),
			strconv.Itoa(r.Unique),
			"https://github.com/" + r.Owner + "/" + r.Name + "/graphs/traffic",
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"repo", "views", "day", "unique", "link"})
	table.AppendBulk(data)
	table.Render()
}

// Remaining functions implement sort.Interface

func (t *TrafficTable) Len() int {
	return len(t.repos)
}

func (t *TrafficTable) Less(i, j int) bool {
	a := t.repos[i]
	b := t.repos[j]
	if a.Views == b.Views {
		if a.Today == b.Today {
			return a.Unique > b.Unique
		}
		return a.Today > b.Today
	}
	return a.Views > b.Views
}

func (t *TrafficTable) Swap(i, j int) {
	a := t.repos[i]
	t.repos[i] = t.repos[j]
	t.repos[j] = a
}
