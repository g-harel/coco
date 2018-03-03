package main

import (
	"flag"
	"fmt"
	"regexp"
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
	names := flag.String("names", "", "comma-separated list of owner names (required)")
	exclude := flag.String("exclude", "", "regexp exclusion mask called on each repo \"owner/repo\"")
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

	if *exclude != "" {
		pattern := regexp.MustCompile(*exclude)
		repos = filter(repos, func(r *Repo) bool {
			return !pattern.Match([]byte(r.Owner + "/" + r.Name))
		})
	}

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
