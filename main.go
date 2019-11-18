package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/g-harel/coco/github"
	"github.com/g-harel/coco/npm"
)

var githubToken = flag.String("github-token", "", "GitHub API token")
var githubUsers = flag.String("github-user", "", "List of GitHub users and orgs whose repos to query (comma separated).")
var npmUsers = flag.String("npm-user", "", "List of NPM users whose packages to query (comma separated).")

func main() {
	flag.Parse()

	var githubTable string
	var npmTable string

	lock := sync.WaitGroup{}
	lock.Add(2)
	go func() {
		users := strings.Split(strings.ReplaceAll(*githubUsers, " ", ""), ",")
		githubTable = github.Repositories(*githubToken, users...)
		lock.Done()
	}()
	go func() {
		users := strings.Split(strings.ReplaceAll(*npmUsers, " ", ""), ",")
		npmTable = npm.Packages(users...)
		lock.Done()
	}()
	lock.Wait()

	fmt.Print(githubTable)
	fmt.Print(npmTable)
}
