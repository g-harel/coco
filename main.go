package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/g-harel/coco/github"
	"github.com/g-harel/coco/npm"
)

func main() {
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

	token, ok := os.LookupEnv("GITHUB_API_TOKEN")
	if !ok {
		fmt.Fprintln(os.Stderr, "Missing GITHUB_API_TOKEN environment variable.")
		os.Exit(1)
	}

	var githubTable string
	var npmTable string

	lock := sync.WaitGroup{}
	lock.Add(2)
	go func() {
		githubTable = github.Repositories(token, users)
		lock.Done()
	}()
	go func() {
		npmTable = npm.Packages("g-harel")
		lock.Done()
	}()
	lock.Wait()

	fmt.Print(githubTable)
	fmt.Print(npmTable)
}
