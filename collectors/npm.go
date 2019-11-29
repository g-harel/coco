package collectors

import (
	"encoding/json"
	"fmt"
	"github.com/g-harel/coco/internal"
	"net/http"
	"net/url"
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

	var wg sync.WaitGroup
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
	u, err := url.Parse(fmt.Sprintf("https://www.npmjs.com/~%v?page=%v", user, page))
	if err != nil {
		return nil, err
	}

	res, err := internal.DefaultLoggingClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user %v: unexpected status code %v", user, res.StatusCode)
	}

	var p npmUserResponse
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func npmFetchPackage(name string) (*npmPackageResponse, error) {
	u, err := url.Parse(fmt.Sprintf("https://www.npmjs.com/package/%v", name))
	if err != nil {
		return nil, err
	}

	res, err := internal.DefaultLoggingClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return nil, err
	}

	var p npmPackageResponse
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
