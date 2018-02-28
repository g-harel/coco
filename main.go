package main

import (
	"flag"
	"fmt"
)

type Repo struct {
	Name   string
	Owner  string
	Views  int
	Today  int
	Unique int
	Error  error
}

func main() {
	names := flag.String("names", "", "comma-separated list of names (required)")
	token := flag.String("token", "", "provide api token (required)")

	flag.Parse()

	if *names == "" || *token == "" {
		fmt.Println("Usage of coco:")
		flag.PrintDefaults()
		return
	}

	client := NewClient(*token, *names)

	repos, err := client.FetchRepos()
	if err != nil {
		panic(err)
	}

	// remove duplicate repos
	visited := make(map[string]bool)
	repos = filter(repos, func(r *Repo) bool {
		if visited[r.Owner+r.Name] {
			return false
		}
		visited[r.Owner+r.Name] = true
		return true
	})

	repos, err = client.FetchTraffic(repos)
	if err != nil {
		panic(err)
	}

	// only show repos with views
	repos = filter(repos, func(r *Repo) bool {
		return r.Views > 0
	})

	NewTrafficTable(repos).Print()
}

func filter(repos []*Repo, f func(*Repo) bool) []*Repo {
	rs := []*Repo{}
	for _, r := range repos {
		if f(r) {
			rs = append(rs, r)
		}
	}
	return rs
}
