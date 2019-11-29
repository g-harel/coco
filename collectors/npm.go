package collectors

import (
	"fmt"
	"github.com/g-harel/coco/internal"
	"net/http"
	"sync"
)

type NpmPackage struct {
	Name   string
	Weekly int
	Total  int
}

type NpmPackageHandler func(*NpmPackage, error)

func NpmPackages(f NpmPackageHandler, users ...string) {
	internal.ExecParallel(len(users), func(i int) {
		npmHandleUser(npmConverterFunc(f), users[i])
	})
}

type npmUserResponse struct {
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

func npmHandleUser(f npmPackageResponseHandler, user string) {
	firstPage, err := npmFetchUser(user, 0)
	if err != nil {
		f(nil, err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		npmHandleUserResponse(f, firstPage)
		wg.Done()
	}()

	remainingPages := firstPage.Packages.Total / firstPage.Pagination.PerPage
	internal.ExecParallel(remainingPages, func(n int) {
		nthPage, err := npmFetchUser(user, n+1)
		if err != nil {
			f(nil, fmt.Errorf("failed to fetch page %v: %v", n, err))
		} else {
			npmHandleUserResponse(f, nthPage)
		}
	})

	wg.Wait()
}

func npmHandleUserResponse(f npmPackageResponseHandler, r *npmUserResponse) {
	internal.ExecParallel(len(r.Packages.Objects), func(n int) {
		f(npmFetchPackage(r.Packages.Objects[n].Name))
	})
}

func npmFetchUser(user string, page int) (*npmUserResponse, error) {
	res := &npmUserResponse{}
	err := internal.HTTPGet(
		fmt.Sprintf("https://www.npmjs.com/~%v?page=%v", user, page),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch user %v page %v: %v", user, page, err)
	}
	return res, nil
}

func npmFetchPackage(name string) (*npmPackageResponse, error) {
	res := &npmPackageResponse{}
	err := internal.HTTPGet(
		fmt.Sprintf("https://www.npmjs.com/package/%v", name),
		http.Header{"x-spiferack": []string{"1"}},
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch user package %v: %v", name, err)
	}
	return res, nil
}
