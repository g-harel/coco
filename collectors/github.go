package collectors

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func githubParsePaginationHeader(h *http.Header) (last int, err error) {
	header := strings.Split(h.Get("Link"), ",")
	for i := 0; i < len(header); i++ {
		if strings.HasSuffix(header[i], `>; rel="last"`) {
			rawURL := strings.TrimSpace(header[i])
			rawURL = strings.TrimPrefix(rawURL, "<")
			rawURL = strings.TrimSuffix(rawURL, `>; rel="last"`)
			url, err := url.Parse(rawURL)
			if err != nil {
				return 0, err
			}
			rawPage := url.Query().Get("page")
			page, err := strconv.Atoi(rawPage)
			if err != nil {
				return 0, err
			}
			return page, nil

		}
	}
	return 0, fmt.Errorf("pagination header value not found")
}

func GithubFetchOwnerRepos(owner string, page int) (*githubRepoListResponse, int, error) {
	res := &githubRepoListResponse{}
	h, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=%v", owner, page),
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", "TODO")}},
		res,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("fetch owner %v page %v: %v", owner, page, err)
	}
	lastPage, err := githubParsePaginationHeader(h)
	if err != nil {
		return nil, 0, fmt.Errorf("parse pagination header: %v", err)
	}
	return res, lastPage, nil
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
