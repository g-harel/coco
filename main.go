package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/github"
	"github.com/g-harel/coco/internal"
)

var githubToken = flag.String("github-token", "", "GitHub API token")
var githubUsers = flag.String("github-user", "", "List of GitHub users and orgs whose repos to query (comma separated).")
var npmUsers = flag.String("npm-user", "", "List of NPM users whose packages to query (comma separated).")

func main() {
	flag.Parse()

	var githubTable string
	var npmTable internal.Table

	lock := sync.WaitGroup{}
	lock.Add(2)
	go func() {
		users := strings.Split(strings.ReplaceAll(*githubUsers, " ", ""), ",")
		githubTable = github.Repositories(*githubToken, users...)
		lock.Done()
	}()
	go func() {
		users := strings.Split(strings.ReplaceAll(*npmUsers, " ", ""), ",")
		npmTable = npmPackages(users...)
		lock.Done()
	}()
	lock.Wait()

	fmt.Print(githubTable)
	fmt.Print(npmTable.String())
}

func npmPackages(users ...string) internal.Table {
	var t internal.Table
	t.Headers("PACKAGE", "DOWNLOADS", "TOTAL", "LINK")
	collectors.NpmPackages(func(p *collectors.NpmPackage, err error) {
		if err != nil {
			internal.LogError("%v\n", err)
			return
		}
		if p.Weekly < 12 {
			return
		}
		link := "https://npmjs.com/package/" + p.Name
		t.Add(p.Name, p.Weekly, p.Total, link)

	}, users...)
	t.Sort(1, 2)
	return t
}
