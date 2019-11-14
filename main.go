package main

import (
	"fmt"
	"os"
	"strings"

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

	fmt.Print(githubRepositories(token, users).String())
	fmt.Print(npm.Packages("g-harel"))
}

func githubRepositories(token string, users []string) github.Repositories {
	gh := github.NewClient(token)

	repos, err := gh.Repositories(users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch repositories: %v", err)
		os.Exit(1)
	}

	// Remove duplicate repositories (usernames might have overlap).
	visited := make(map[string]bool)
	repos = repos.Filter(func(r *github.Repository) bool {
		if visited[r.Owner+r.Name] {
			return false
		}
		visited[r.Owner+r.Name] = true
		return true
	})

	// Fetch traffic data for all repositories.
	repos = gh.Traffic(repos)

	// Remove repos with errors or no reported views (in the past two weeks).
	repos = repos.Filter(func(r *github.Repository) bool {
		if r.Error != nil {
			// Fetching errors are swallowed to avoid crowding the output (subject to change).
			return false
		}
		return r.Views != 0
	})

	return repos
}
