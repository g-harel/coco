package collectors

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

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

func GithubRepos(f GithubRepoHandler, token string, owners ...string) {
	internal.ExecParallel(len(owners), func(i int) {
		githubHandleOwner(githubConverterFunc(f), token, owners[i])
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

	Name  string
	Owner string
}

type githubRepoResponseHandler func(*githubRepoViewsResponse, error)

func githubHandleOwner(f githubRepoResponseHandler, token, owner string) {
	firstPage, last, err := githubFetchOwnerRepos(token, owner, 0)
	if err != nil {
		f(nil, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		githubHandleRepoListResponse(f, token, firstPage)
		wg.Done()
	}()

	remainingPages := last - 1
	internal.ExecParallel(remainingPages, func(n int) {
		nthPage, _, err := githubFetchOwnerRepos(token, owner, n+1)
		if err != nil {
			f(nil, fmt.Errorf("fetch page %v: %v", n, err))
		} else {
			githubHandleRepoListResponse(f, token, nthPage)
		}
	})

	wg.Wait()
}

func githubHandleRepoListResponse(f githubRepoResponseHandler, token string, r githubRepoListResponse) {
	internal.ExecParallel(len(r), func(n int) {
		f(githubFetchRepoViews(token, r[n].Owner.Login, r[n].Name))
	})
}

func githubConverterFunc(f GithubRepoHandler) githubRepoResponseHandler {
	return func(r *githubRepoViewsResponse, err error) {
		if err != nil {
			f(nil, err)
			return
		}
		p := &GithubRepo{
			Owner:  r.Owner,
			Name:   r.Name,
			Views:  r.Count,
			Unique: r.Uniques,
		}
		// TODO calc today
		f(p, nil)
	}
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

func githubFetchOwnerRepos(token, owner string, page int) (githubRepoListResponse, int, error) {
	res := githubRepoListResponse{}
	h, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=%v", owner, page),
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}},
		&res,
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

func githubFetchRepoViews(token, owner, name string) (*githubRepoViewsResponse, error) {
	res := &githubRepoViewsResponse{}
	_, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch repo views %v: %v", name, err)
	}
	res.Owner = owner
	res.Name = name
	return res, nil
}
