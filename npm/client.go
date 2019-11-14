package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type UserPage struct {
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

type PackageData struct {
	Package   string `json:"package"`
	Downloads []struct {
		Downloads int `json:"downloads"`
	} `json:"downloads"`
}

func fetchUserPage(user string, page int) (*UserPage, error) {
	u, err := url.Parse(fmt.Sprintf("https://www.npmjs.com/~%v?page=%v", user, page))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return nil, err
	}

	var p UserPage
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func fetchAllUserPackageNames(user string) ([]string, error) {
	firstPage, err := fetchUserPage(user, 0)
	if err != nil {
		return nil, err
	}

	totalPages := firstPage.Packages.Total/firstPage.Pagination.PerPage + 1
	pages := make([]*UserPage, totalPages)
	pages[0] = firstPage

	var wg sync.WaitGroup
	for i := 1; i < totalPages; i++ {
		wg.Add(1)
		go func(page int) {
			nthPage, err := fetchUserPage(user, page)
			if err != nil {
				fmt.Printf("failed to fetch page %v: %v", page, err)
			}
			pages[page] = nthPage
			wg.Done()
		}(i)
	}
	wg.Wait()

	packages := []string{}
	for i := 0; i < len(pages); i++ {
		page := pages[i]
		for j := 0; j < len(page.Packages.Objects); j++ {
			packages = append(packages, page.Packages.Objects[j].Name)
		}
	}

	return packages, nil
}

func fetchPackageData(name string) (*PackageData, error) {
	u, err := url.Parse(fmt.Sprintf("https://www.npmjs.com/package/%v", name))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return nil, err
	}

	var p PackageData
	err = json.NewDecoder(res.Body).Decode(&p)

	return &p, nil
}

func convert(d *PackageData) *Pkg {
	p := &Pkg{
		Name:   d.Package,
		Weekly: 0,
		Total:  0,
	}
	if len(d.Downloads) > 0 {
		p.Weekly = d.Downloads[len(d.Downloads)-1].Downloads
		for i := 0; i < len(d.Downloads); i++ {
			p.Total += d.Downloads[i].Downloads
		}
	}
	return p
}

func Packages(user string) string {
	p, err := fetchAllUserPackageNames(user)
	if err != nil {
		panic(err)
	}

	packages := make(PackageList, len(p))
	var wg sync.WaitGroup
	visited := map[string]bool{}
	for i := 0; i < len(p); i++ {
		visited[p[i]] = true
		wg.Add(1)
		go func(name string, i int) {
			d, err := fetchPackageData(name)
			if err != nil {
				panic(err)
			}
			packages[i] = convert(d)
			wg.Done()
		}(p[i], i)
	}
	wg.Wait()

	return packages.String()
}
