package collectors

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/httpc"
)

type GithubRepo struct {
	Name   string
	Owner  string
	Stars  int
	Views  int
	Today  int
	Unique int
}

type GithubRepoHandler func(*GithubRepo, error)

func GithubRepos(f GithubRepoHandler, token string, owners []string) {
	exec.Parallel(len(owners), func(i int) {
		githubHandleOwner(f, token, owners[i])
	})
}

type githubRepoListResponse []struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Stars int `json:"stargazers_count"`
}

type githubRepoViewsResponse struct {
	Count   int `json:"count"`
	Uniques int `json:"uniques"`
	Views   []struct {
		Timestamp string `json:"timestamp"`
		Count     int    `json:"count"`
	} `json:"views"`
}

type githubRepoResponseHandler func(*githubRepoViewsResponse, error)

func githubHandleOwner(f GithubRepoHandler, token, owner string) {
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
	exec.Parallel(remainingPages, func(n int) {
		nthPageURL, _ := url.Parse(lastURL.String())
		query := nthPageURL.Query()
		query.Del("page")
		query.Add("page", strconv.Itoa(n+2))
		nthPageURL.RawQuery = query.Encode()
		nthPage := githubRepoListResponse{}
		_, err := httpc.Get(
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

func githubHandleRepoListResponse(f GithubRepoHandler, token string, r githubRepoListResponse) {
	exec.Parallel(len(r), func(n int) {
		convertedHandler := githubConverterFunc(f, r[n].Owner.Login, r[n].Name, r[n].Stars)
		convertedHandler(githubFetchRepoViews(token, r[n].Owner.Login, r[n].Name))
	})
}

func githubConverterFunc(f GithubRepoHandler, owner, name string, stars int) githubRepoResponseHandler {
	return func(r *githubRepoViewsResponse, err error) {
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
		p := &GithubRepo{
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
	h, err := httpc.Get(
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
	_, err := httpc.Get(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		githubTokenHeader(token),
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch repo views %v: %v", name, err)
	}
	return res, nil
}

func githubTokenHeader(token string) http.Header {
	return http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}}
}
