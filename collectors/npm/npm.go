package npm

import (
	"fmt"
	"sync"

	"github.com/g-harel/coco/internal/exec"
)

func Packages(f PackageHandler, owners []string) {
	exec.Parallel(len(owners), func(i int) {
		handleOwner(converterFunc(f), owners[i])
	})
}

func handleOwner(f packageResponseHandler, owner string) {
	firstPage, err := fetchOwner(owner, 0)
	if err != nil {
		f(nil, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		handleOwnerResponse(f, firstPage)
		wg.Done()
	}()

	remainingPages := firstPage.Packages.Total / firstPage.Pagination.PerPage
	exec.Parallel(remainingPages, func(n int) {
		nthPage, err := fetchOwner(owner, n+1)
		if err != nil {
			f(nil, fmt.Errorf("fetch page %v: %v", n, err))
		} else {
			handleOwnerResponse(f, nthPage)
		}
	})

	wg.Wait()
}

func handleOwnerResponse(f packageResponseHandler, r *ownerResponse) {
	exec.Parallel(len(r.Packages.Objects), func(n int) {
		f(fetchPackage(r.Packages.Objects[n].Name))
	})
}
