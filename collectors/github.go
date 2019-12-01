package collectors

import (
	"fmt"
	"net/http"

	"github.com/g-harel/coco/internal"
)

type GithubRepo struct {
	Name   string
	Owner  string
	Views  int
	Today  int
	Unique int
}

type GithubRepoHandler func(*GithubRepo, error)

func GithubRepos(f GithubRepoHandler, owners ...string) {
	internal.ExecParallel(len(owners), func(i int) {
		githubHandleOwner(githubConverterFunc(f), owners[i])
	})
}

type githubRepoListResponse []struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
}

type githubRepoViewsResponse struct {
	Count   int `json:"count"`
	Uniques int `json:"uniques"`
	Views   []struct {
		Timestamp string `json:"timestamp"`
		Count     int    `json:"count"`
		Uniques   int    `json:"uniques"`
	} `json:"views"`
}

type githubRepoResponseHandler func(*githubRepoViewsResponse, error)

func githubConverterFunc(f GithubRepoHandler) githubRepoResponseHandler {
	return func(r *githubRepoViewsResponse, err error) {
		if err != nil {
			f(nil, err)
			return
		}
		p := &GithubRepo{
			Views:  r.Count,
			Unique: r.Uniques,
		}
		f(p, nil)
	}
}

func githubHandleOwner(f githubRepoResponseHandler, owner string) {

}

func githubParsePaginationHeader(h http.Header) (last int, ok bool) {
	return 0, false
}

func GithubFetchOwnerRepos(owner string, page int) (*githubRepoListResponse, error) {
	res := &githubRepoListResponse{}
	_, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=%v", owner, page),
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", "TODO")}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner %v page %v: %v", owner, page, err)
	}
	// TODO add pagination
	return res, nil
}

func githubFetchRepoViews(owner, name string) (*githubRepoViewsResponse, error) {
	res := &githubRepoViewsResponse{}
	_, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", "TODO")}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner package %v: %v", name, err)
	}
	return res, nil
}
