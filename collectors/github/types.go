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

type repoListResponse []struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Stars int `json:"stargazers_count"`
}

type repoViewsResponse struct {
	Count   int `json:"count"`
	Uniques int `json:"uniques"`
	Views   []struct {
		Timestamp string `json:"timestamp"`
		Count     int    `json:"count"`
	} `json:"views"`
}

type repoResponseHandler func(*repoViewsResponse, error)

func converterFunc(f RepoHandler, owner, name string, stars int) repoResponseHandler {
	return func(r *repoViewsResponse, err error) {
		if err != nil {
			f(nil, err)
			return
		}
		today := 0
		nowPrefix := time.Now().Format("2006-01-02")
		for i := 0; i < len(r.Views); i++ {
			if strings.HasPrefix(r.Views[i].Timestamp, nowPrefix) {
				today += r.Views[i].Count
			}
		}
		p := &Repo{
			Name:   name,
			Owner:  owner,
			Stars:  stars,
			Views:  r.Count,
			Today:  today,
			Unique: r.Uniques,
		}
		f(p, nil)
	}
}
