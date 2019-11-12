package npm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func Packages(user string) ([]string, error) {
	firstPage, err := fetchUserPage(user, 0)
	if err != nil {
		return nil, err
	}

	totalPages := firstPage.Packages.Total/firstPage.Pagination.PerPage + 1
	pages := []*UserPage{firstPage}

	for i := 1; i < totalPages; i++ {
		nthPage, err := fetchUserPage(user, i)
		if err != nil {
			fmt.Printf("failed to fetch page %v: %v", i, err)
		}
		pages = append(pages, nthPage)
	}

	packages := []string{}
	for i := 0; i < len(pages); i++ {
		page := pages[i]
		for j := 0; j < len(page.Packages.Objects); j++ {
			packages = append(packages, page.Packages.Objects[j].Name)
		}
	}

	return packages, nil
}

func Package(name string) (*PackageData, error) {
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
