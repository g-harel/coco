package github

import (
	"fmt"

	"github.com/g-harel/coco/internal/exec"
)

// HandleOwner fetches all an owner's repos and calls the
// next handler with them. Using the first request, it reads
// all paginated data in parallel.
func handleOwner(f repoHandler, token, owner string) {
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

// HandleReposResponse reads all repos from a single page
// and fetches its views data.
func handleReposResponse(f repoHandler, token string, l reposResponse) {
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
