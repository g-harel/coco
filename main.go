package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/internal"
)

var githubToken = flag.String("github-token", "", "GitHub API token")
var githubOwners = flag.String("github-owner", "", "List of GitHub owners whose repos to query (comma separated).")
var githubViews = flag.Int("github-views", 0, "Show repos if they have this quantity of views.")
var githubToday = flag.Int("github-today", 0, "Show repos if they have this quantity of views today.")

var npmOwners = flag.String("npm-owner", "", "List of NPM owners whose packages to query (comma separated).")
var npmWeekly = flag.Int("npm-weekly", 0, "Show repos if they have this quantity of weekly downloads.")

func main() {
	flag.Parse()

	githubTable := internal.Table{}
	npmTable := internal.Table{}

	lock := sync.WaitGroup{}
	lock.Add(2)
	go func() {
		owners := strings.Split(strings.ReplaceAll(*githubOwners, " ", ""), ",")
		githubTable = collectGithubPackages(*githubToken, owners)
		lock.Done()
	}()
	go func() {
		owners := strings.Split(strings.ReplaceAll(*npmOwners, " ", ""), ",")
		npmTable = collectNpmPackages(owners)
		lock.Done()
	}()
	lock.Wait()

	fmt.Print(githubTable.String())
	fmt.Print(npmTable.String())
}

func collectNpmPackages(owners []string) internal.Table {
	t := internal.Table{}
	t.Headers(
		"PACKAGE",
		"DOWNLOADS",
		"TOTAL",
		"LINK",
	)
	collectors.NpmPackages(func(p *collectors.NpmPackage, err error) {
		if err != nil {
			internal.LogError("%v\n", err)
			return
		}
		if p.Weekly < *npmWeekly {
			return
		}
		link := "https://npmjs.com/package/" + p.Name
		t.Add(
			p.Name,
			p.Weekly,
			p.Total,
			link,
		)

	}, owners)
	t.Sort(1, 2)
	return t
}

func collectGithubPackages(token string, owners []string) internal.Table {
	t := internal.Table{}
	t.Headers(
		"REPO",
		"VIEWS",
		"UNIQUE",
		"TODAY",
		"LINK",
	)
	collectors.GithubRepos(func(r *collectors.GithubRepo, err error) {
		if err != nil {
			internal.LogError("%v\n", err)
			return
		}
		if r.Today < *githubToday && r.Views < *githubViews {
			return
		}
		t.Add(
			fmt.Sprintf("%v*%v", r.Name, r.Stars),
			r.Views,
			r.Unique,
			r.Today,
			fmt.Sprintf("https://github.com/%v/%v/graphs/traffic", r.Owner, r.Name),
		)
	}, token, owners)
	t.Sort(1, 3, 2)
	return t
}
