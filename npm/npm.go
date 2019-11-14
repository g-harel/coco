package npm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
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
	sort.Sort(packages)

	data := [][]string{}
	for _, p := range packages {
		if p.Weekly < 12 {
			continue
		}
		data = append(data, []string{
			p.Name,
			strconv.Itoa(p.Weekly),
			strconv.Itoa(p.Total),
			"https://npmjs.com/package/" + p.Name,
		})
	}

	buf := &bytes.Buffer{}
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"package", "downloads", "total", "link"})
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_LEFT,
	})
	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (packages packageStatsList) Len() int {
	return len(packages)
}

func (packages packageStatsList) Less(i, j int) bool {
	a := packages[i]
	b := packages[j]
	// Sort by: weekly downloads -> total downloads -> name
	if a.Weekly == b.Weekly {
		if a.Total == b.Total {
			return strings.Compare(a.Name, b.Name) < 0
		}
		return a.Total > b.Total
	}
	return a.Weekly > b.Weekly
}

func (packages packageStatsList) Swap(i, j int) {
	packages[i], packages[j] = packages[j], packages[i]
}

func fetchUserPage(user string, page int) (*userPage, error) {
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

	totalPages := firstPage.Packages.Total/firstPage.Pagination.PerPage + 1
	pages := make([]*userPage, totalPages)
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

func fetchPackageData(name string) (*packageData, error) {
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

func Packages(user string) string {
	p, err := fetchAllUserPackageNames(user)
	if err != nil {
		panic(err)
	}

	packages := make(packageStatsList, len(p))
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