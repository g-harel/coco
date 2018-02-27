package main

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	name   string
}

func NewClient(name, token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &Client{
		github.NewClient(tc),
		name,
	}
}

func (c *Client) FetchRepos() ([]*Repo, error) {
	list, _, err := c.client.Repositories.List(
		context.Background(),
		c.name,
		&github.RepositoryListOptions{
			Type: "all",
		},
	)
	if err != nil {
		return nil, err
	}

	repos := make([]*Repo, len(list))
	for i, l := range list {
		repos[i] = &Repo{
			Owner: l.GetOwner().GetLogin(),
			Name:  l.GetName(),
		}
	}

	return repos, nil
}

func (c *Client) FetchTraffic(repos []*Repo) ([]*Repo, error) {
	rch := make(chan *Repo, len(repos))

	for i := 0; i < len(repos); i++ {
		go func(i int) {
			v, _, err := c.client.Repositories.ListTrafficViews(
				context.Background(),
				repos[i].Owner,
				repos[i].Name,
				nil,
			)
			if err != nil {
				rch <- &Repo{
					Error: err,
				}
				return
			}

			today := 0
			for _, td := range v.Views {
				if isToday(td.Timestamp) {
					today++
				}
			}

			rch <- &Repo{
				Name:   repos[i].Name,
				Owner:  repos[i].Owner,
				Views:  v.GetCount(),
				Today:  today,
				Unique: v.GetUniques(),
			}
		}(i)
	}

	var err error
	r := make([]*Repo, len(repos))
	for i := 0; i < len(repos); i++ {
		r[i] = <-rch
		if r[i].Error != nil {
			err = r[i].Error
			break
		}
	}

	return r, err
}

func isToday(t *github.Timestamp) bool {
	y, m, d := time.Now().Date()
	cy, cm, cd := t.Date()
	return y == cy && m == cm && d == cd
}
