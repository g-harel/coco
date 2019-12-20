package npm

import (
	"fmt"

	"github.com/g-harel/coco/internal/exec"
)

// HandleOwner fetches all an owner's packages and calls the
// next handler with them. Using the first request, it reads
// all paginated data in parallel.
func handleOwner(f pkgHandler, owner string) {
	firstPage, err := fetchOwner(owner, 0)
	if err != nil {
		f(nil, err)
		return
	}

	exec.Parallel(
		func() {
			handleOwnerResponse(f, firstPage)
		},
		func() {
			remainingPages := firstPage.Packages.Total / firstPage.Pagination.PerPage
			exec.ParallelN(remainingPages, func(n int) {
				// Skip already loaded first page.
				nthPage, err := fetchOwner(owner, n+1)
				if err != nil {
					f(nil, fmt.Errorf("fetch page %v: %v", n, err))
				} else {
					handleOwnerResponse(f, nthPage)
				}
			})
		},
	)
}

// HandleOwnerResponse reads all packages from a single page
// of owner data and fetches package download data.
func handleOwnerResponse(f pkgHandler, r *ownerResponse) {
	exec.ParallelN(len(r.Packages.Objects), func(n int) {
		r, err := fetchPackage(r.Packages.Objects[n].Name)
		if err != nil {
			f(nil, err)
		} else {
			f(convert(r), nil)
		}
	})
}
