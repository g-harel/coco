package github

import (
	"strings"
	"time"
)

type Repo struct {
	Name   string
	Owner  string
	Stars  int
	Views  int
	Today  int
	Unique int
}

type RepoHandler func(*Repo, error)

type reposResponse []struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Stars int `json:"stargazers_count"`
}

type viewsResponse struct {
	Count   int `json:"count"`
	Uniques int `json:"uniques"`
	Views   []struct {
		Timestamp string `json:"timestamp"`
		Count     int    `json:"count"`
	} `json:"views"`
}

func convert(v *viewsResponse) *Repo {
	today := 0
	nowPrefix := time.Now().Format("2006-01-02")
	for i := 0; i < len(v.Views); i++ {
		if strings.HasPrefix(v.Views[i].Timestamp, nowPrefix) {
			today += v.Views[i].Count
		}
	}
	r := &Repo{
		Views:  v.Count,
		Today:  today,
		Unique: v.Uniques,
	}
	return r
}
