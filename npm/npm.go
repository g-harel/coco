package npm

import (
	"encoding/json"
	"fmt"
	"github.com/g-harel/coco/internal"
	"net/http"
	"net/url"
)

type userPage struct {
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

type packageData struct {
	Package   string `json:"package"`
	Downloads []struct {
		Downloads int `json:"downloads"`
	} `json:"downloads"`
}

type packageStats struct {
	Name   string
	Weekly int
	Total  int
}

type packageStatsList []*packageStats

func (packages packageStatsList) String() string {
	t := internal.Table{}
	t.Headers("PACKAGE", "DOWNLOADS", "TOTAL", "LINK")
	for _, p := range packages {
		if p.Weekly < 12 {
			continue
		}
		t.Add(
			p.Name,
			p.Weekly,
			p.Total,
			"https://npmjs.com/package/"+p.Name,
		)
	}
	return t.Format(1, 2)
}

func fetchUserPage(user string, page int) (*userPage, error) {
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

	var p userPage
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

	if firstPage.Packages.Total == 0 {
		return []string{}, nil
	}

	totalPages := firstPage.Packages.Total/firstPage.Pagination.PerPage + 1
	pages := make([]*userPage, totalPages)

	internal.ExecParallel(totalPages, func(page int) {
		if page == 0 {
			pages[page] = firstPage
			return
		}
		nthPage, err := fetchUserPage(user, page)
		if err != nil {
			fmt.Printf("failed to fetch page %v: %v", page, err)
		}
		pages[page] = nthPage
	})

	packages := []string{}
	for i := 0; i < len(pages); i++ {
		page := pages[i]
		for j := 0; j < len(page.Packages.Objects); j++ {
			packages = append(packages, page.Packages.Objects[j].Name)
		}
	}

	return packages, nil
}

func fetchPackageData(name string) (*packageData, error) {
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

	var p packageData
	err = json.NewDecoder(res.Body).Decode(&p)

	return &p, nil
}

func convert(d *packageData) *packageStats {
	p := &packageStats{
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

func Packages(users ...string) string {
	packageNames := []string{}
	internal.ExecParallel(len(users), func(i int) {
		p, err := fetchAllUserPackageNames(users[i])
		if err != nil {
			internal.LogError("fetch user package names: %v", err)
			return
		}
		internal.ExecSafe(func() {
			packageNames = append(packageNames, p...)
		})
	})

	packageStats := make(packageStatsList, len(packageNames))
	internal.ExecParallel(len(packageNames), func(i int) {
		data, err := fetchPackageData(packageNames[i])
		if err != nil {
			internal.LogError("fetch package data: %v", err)
			return
		}
		packageStats[i] = convert(data)
	})

	return packageStats.String()
}
