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
	firstPage, lastURL, err := githubFetchInitialOwnerRepo(token, owner)
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

	if lastURL == nil {
		wg.Wait()
		return
	}

	last, err := strconv.Atoi(lastURL.Query().Get("page"))
	if err != nil {
		f(nil, fmt.Errorf("parse last pagination index: %v", err))
		return
	}

	remainingPages := last - 1
	internal.ExecParallel(remainingPages, func(n int) {
		nthPageURL, _ := url.Parse(lastURL.String())
		query := nthPageURL.Query()
		query.Del("page")
		query.Add("page", strconv.Itoa(n+2))
		nthPageURL.RawQuery = query.Encode()
		nthPage := githubRepoListResponse{}
		_, err := internal.HTTPGet(
			nthPageURL.String(),
			githubTokenHeader(token),
			&nthPage,
		)
		if err != nil {
			f(nil, fmt.Errorf("fetch owner %v page %v: %v", owner, n+1, err))
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

func githubParsePaginationHeaderLastURL(h *http.Header) (*url.URL, error) {
	header := strings.Split(h.Get("Link"), ",")
	for i := 0; i < len(header); i++ {
		if strings.HasSuffix(header[i], `>; rel="last"`) {
			rawURL := strings.TrimSpace(header[i])
			rawURL = strings.TrimPrefix(rawURL, "<")
			rawURL = strings.TrimSuffix(rawURL, `>; rel="last"`)
			url, err := url.Parse(rawURL)
			if err != nil {
				return nil, fmt.Errorf("parse pagination url: %v", err)
			}
			return url, nil

		}
	}
	return nil, fmt.Errorf("pagination header value not found")
}

func githubFetchInitialOwnerRepo(token, owner string) (githubRepoListResponse, *url.URL, error) {
	res := githubRepoListResponse{}
	h, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=1", owner),
		githubTokenHeader(token),
		&res,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch owner initial page %v: %v", owner, err)
	}
	lastPageURL, err := githubParsePaginationHeaderLastURL(h)
	if err != nil {
		return res, nil, nil
	}
	return res, lastPageURL, nil

}

func githubFetchRepoViews(token, owner, name string) (*githubRepoViewsResponse, error) {
	res := &githubRepoViewsResponse{}
	_, err := internal.HTTPGet(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		githubTokenHeader(token),
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch repo views %v: %v", name, err)
	}
	res.Owner = owner
	res.Name = name
	return res, nil
}

func githubTokenHeader(token string) http.Header {
	return http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}}
}
