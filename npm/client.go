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
	u, err := url.Parse(fmt.Sprintf("https://www.npmjs.com/~%v?page=0", user))
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
	total := 1
	page := 0
	perPage := 0
	packages := []string{}
	for (page+1)*perPage < total {
		p, err := fetchUserPage(user, page)
		if err != nil {
			return nil, err
		}
		total = p.Packages.Total
		page = p.Pagination.Page
		perPage = p.Pagination.PerPage
		for i := 0; i < len(p.Packages.Objects); i++ {
			packages = append(packages, p.Packages.Objects[i].Name)
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
