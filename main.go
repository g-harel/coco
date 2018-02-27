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
	name := flag.String("name", "", "specify username (required)")
	token := flag.String("token", "", "provide api token (required)")

	flag.Parse()

	if *name == "" || *token == "" {
		fmt.Println("Usage of coco:")
		flag.PrintDefaults()
		return
	}

	client := NewClient(*name, *token)

	repos, err := client.FetchRepos()
	if err != nil {
		panic(err)
	}

	repos, err = client.FetchTraffic(repos)
	if err != nil {
		panic(err)
	}

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
