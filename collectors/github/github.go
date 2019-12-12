package github

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/httpc"
)

func Repos(f RepoHandler, token string, owners []string) {
	exec.Parallel(len(owners), func(i int) {
		handleOwner(f, token, owners[i])
	})
}

func handleOwner(f RepoHandler, token, owner string) {
	firstPage, lastURL, err := fetchInitialOwnerRepo(token, owner)
	if err != nil {
		f(nil, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		handleRepoListResponse(f, token, firstPage)
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
		nthPage := repoListResponse{}
		_, err := httpc.Get(
			nthPageURL.String(),
			tokenHeader(token),
			&nthPage,
		)
		if err != nil {
			f(nil, fmt.Errorf("fetch owner %v page %v: %v", owner, n+1, err))
		} else {
			handleRepoListResponse(f, token, nthPage)
		}
	})

	wg.Wait()
}

func handleRepoListResponse(f RepoHandler, token string, r repoListResponse) {
	exec.Parallel(len(r), func(n int) {
		convertedHandler := converterFunc(f, r[n].Owner.Login, r[n].Name, r[n].Stars)
		convertedHandler(fetchRepoViews(token, r[n].Owner.Login, r[n].Name))
	})
}
