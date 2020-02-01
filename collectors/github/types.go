package github

import (
	"strings"
	"time"
)

// Repo is extracted views data.
type repo struct {
	Name   string
	Owner  string
	Stars  int
	Views  int
	Today  int
	Unique int
}

// RepoHandler accepts and handles repo views data.
type repoHandler func(*repo, error)

// ReposResponse represents the response data for a request
// for an owner's repositories.
type reposResponse []struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Stars int `json:"stargazers_count"`
}

// ViewsResponse represents the response data for a request
// for repo views data.
type viewsResponse struct {
	Count   int `json:"count"`
	Uniques int `json:"uniques"`
	Views   []struct {
		Timestamp string `json:"timestamp"`
		Count     int    `json:"count"`
	} `json:"views"`
}

// UserResponse represents the response data from a request
// for user information.
type userResponse struct {
	Login     string `json:"login"`
	Followers int    `json:"followers"`
}

// Convert converts between the HTTP response and extracted
// repo views data.
func convert(v *viewsResponse) *repo {
	today := 0
	nowPrefix := time.Now().Format("2006-01-02")
	for i := 0; i < len(v.Views); i++ {
		if strings.HasPrefix(v.Views[i].Timestamp, nowPrefix) {
			today += v.Views[i].Count
		}
	}
	r := &repo{
		Views:  v.Count,
		Today:  today,
		Unique: v.Uniques,
	}
	return r
}
