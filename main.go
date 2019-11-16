package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/g-harel/coco/github"
	"github.com/g-harel/coco/npm"
)

var githubToken = flag.String("github-token", "", "GitHub API token")
var githubUser = flag.String("github-user", "", "List of GitHub users and orgs whose repos to query (comma separated).")
var npmUser = flag.String("npm-user", "", "List of NPM users whose packages to query (comma separated).")

func main() {
	flag.Parse()

	println(*githubToken)
	println(*githubUser)
	println(*npmUser)

	users := os.Args[1:]
	if len(users) == 0 {
		fmt.Println(strings.TrimSpace(`
Usage: coco [USER]...
List repository traffic for USER(s).

Examples:
  coco username
  coco orgname username

Looks for GitHub api key in GITHUB_API_TOKEN environment variable.
Traffic can only be collected from repositories that your account has push access to.
		`))
		return
	}

	var githubTable string
	var npmTable string

	lock := sync.WaitGroup{}
	lock.Add(2)
	go func() {
		githubTable = github.Repositories(*githubToken, []string{*githubUser})
		lock.Done()
	}()
	go func() {
		npmTable = npm.Packages(*npmUser)
		lock.Done()
	}()
	lock.Wait()

	fmt.Print(githubTable)
	fmt.Print(npmTable)
}
