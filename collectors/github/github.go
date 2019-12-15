package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/g-harel/coco/internal/exec"
)

func Repos(f RepoHandler, token string, owners []string) {
	exec.ParallelN(len(owners), func(n int) {
		handleOwner(f, token, owners[n])
	})
}

func handleOwner(f RepoHandler, token, owner string) {
	firstPage, responseHeaders, err := fetchFirstRepos(token, owner)
	if err != nil {
		f(nil, err)
		return
	}

	exec.Parallel(
		func() {
			handleReposResponse(f, token, firstPage)
		},
		func() {
			pages, err := generatePaginatedURLs(responseHeaders)
			if err != nil {
				f(nil, err)
				return
			}
			exec.ParallelN(len(pages)-1, func(n int) {
				nthPage := reposResponse{}
				_, err := fetchGeneric(pages[n+1], token, &nthPage)
				if err != nil {
					f(nil, fmt.Errorf("fetch owner %v page %v: %v", owner, n+1, err))
				} else {
					handleReposResponse(f, token, nthPage)
				}
			})
		},
	)
}

func handleReposResponse(f RepoHandler, token string, l reposResponse) {
	exec.ParallelN(len(l), func(n int) {
		v, err := fetchViews(token, l[n].Owner.Login, l[n].Name)
		if err != nil {
			f(nil, err)
		} else {
			r := convert(v)
			r.Owner = l[n].Owner.Login
			r.Name = l[n].Name
			r.Stars = l[n].Stars
			f(r, nil)
		}
	})
}

func generatePaginatedURLs(h *http.Header) ([]string, error) {
	links := strings.Split(h.Get("Link"), ",")
	rawLastURL := ""
	for i := 0; i < len(links); i++ {
		if strings.HasSuffix(links[i], `>; rel="last"`) {
			rawURL := strings.TrimSpace(links[i])
			rawURL = strings.TrimPrefix(rawURL, "<")
			rawLastURL = strings.TrimSuffix(rawURL, `>; rel="last"`)
		}
	}
	if rawLastURL == "" {
		// Pagination header is missing when the response contains all data.
		return []string{rawLastURL}, nil
	}
	lastURL, err := url.Parse(rawLastURL)
	if err != nil {
		return nil, fmt.Errorf("parse pagination url: %v", err)
	}
	lastPageIndex, err := strconv.Atoi(lastURL.Query().Get("page"))
	if err != nil {
		return nil, fmt.Errorf("parse last pagination index: %v", err)
	}
	urls := []string{}
	for i := 1; i <= lastPageIndex; i++ {
		nthPageURL, _ := url.Parse(lastURL.String())
		query := nthPageURL.Query()
		query.Del("page")
		query.Add("page", strconv.Itoa(i))
		nthPageURL.RawQuery = query.Encode()
		urls = append(urls, nthPageURL.String())
	}
	return urls, nil
}
