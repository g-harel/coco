package collectors

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/g-harel/coco/internal/httpc"
	"github.com/g-harel/coco/internal/exec"
)

type NpmPackage struct {
	Name   string
	Weekly int
	Total  int
}

type NpmPackageHandler func(*NpmPackage, error)

func NpmPackages(f NpmPackageHandler, owners []string) {
	exec.Parallel(len(owners), func(i int) {
		npmHandleOwner(npmConverterFunc(f), owners[i])
	})
}

type npmOwnerResponse struct {
	Packages struct {
		Total   int `json:"total"`
		Objects []struct {
			Name string `json:"name"`
		} `json:"objects"`
	} `json:"packages"`
	Pagination struct {
		PerPage int `json:"perPage"`
		Page    int `json:"page"`
	}
}

type npmPackageResponse struct {
	Package   string `json:"package"`
	Downloads []struct {
		Downloads int `json:"downloads"`
	} `json:"downloads"`
}

type npmPackageResponseHandler func(*npmPackageResponse, error)

func npmHandleOwner(f npmPackageResponseHandler, owner string) {
	firstPage, err := npmFetchOwner(owner, 0)
	if err != nil {
		f(nil, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		npmHandleOwnerResponse(f, firstPage)
		wg.Done()
	}()

	remainingPages := firstPage.Packages.Total / firstPage.Pagination.PerPage
	exec.Parallel(remainingPages, func(n int) {
		nthPage, err := npmFetchOwner(owner, n+1)
		if err != nil {
			f(nil, fmt.Errorf("fetch page %v: %v", n, err))
		} else {
			npmHandleOwnerResponse(f, nthPage)
		}
	})

	wg.Wait()
}

func npmHandleOwnerResponse(f npmPackageResponseHandler, r *npmOwnerResponse) {
	exec.Parallel(len(r.Packages.Objects), func(n int) {
		f(npmFetchPackage(r.Packages.Objects[n].Name))
	})
}

func npmConverterFunc(f NpmPackageHandler) npmPackageResponseHandler {
	return func(r *npmPackageResponse, err error) {
		if err != nil {
			f(nil, err)
			return
		}
		p := &NpmPackage{
			Name:   r.Package,
			Weekly: 0,
			Total:  0,
		}
		if len(r.Downloads) > 0 {
			p.Weekly = r.Downloads[len(r.Downloads)-1].Downloads
			for i := 0; i < len(r.Downloads); i++ {
				p.Total += r.Downloads[i].Downloads
			}
		}
		f(p, nil)
	}
}

func npmFetchOwner(owner string, page int) (*npmOwnerResponse, error) {
	res := &npmOwnerResponse{}
	_, err := httpc.Get(
		fmt.Sprintf("https://www.npmjs.com/~%v?page=%v", owner, page),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner %v page %v: %v", owner, page, err)
	}
	return res, nil
}

func npmFetchPackage(name string) (*npmPackageResponse, error) {
	res := &npmPackageResponse{}
	_, err := httpc.Get(
		fmt.Sprintf("https://www.npmjs.com/package/%v", name),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch owner package %v: %v", name, err)
	}
	return res, nil
}
