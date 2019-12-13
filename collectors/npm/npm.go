package npm

import (
	"fmt"

	"github.com/g-harel/coco/internal/exec"
)

func Packages(f PackageHandler, owners []string) {
	exec.ParallelN(len(owners), func(n int) {
		handleOwner(f, owners[n])
	})
}

func handleOwner(f PackageHandler, owner string) {
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

func handleOwnerResponse(f PackageHandler, r *ownerResponse) {
	exec.ParallelN(len(r.Packages.Objects), func(n int) {
		r, err := fetchPackage(r.Packages.Objects[n].Name)
		if err != nil {
			f(nil, err)
		} else {
			f(convert(r), nil)
		}
	})
}
