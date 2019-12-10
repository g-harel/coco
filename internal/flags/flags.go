package flags

import (
	"flag"
)

var (
	LogInfo   = flag.Bool("log-info", false, "Log miscellaneous info.")
	LogErrors = flag.Bool("log-error", true, "Log errors.")

	GithubToken  = flag.String("github-token", "", "GitHub API token.")
	GithubOwners = flag.String("github-owner", "", "List of GitHub owners whose repos to query (comma separated).")
	GithubViews  = flag.Int("github-views", 1, "Show repos if they have this quantity of views.")
	GithubToday  = flag.Int("github-today", 1, "Show repos if they have this quantity of views today.")
	GithubStars  = flag.Int("github-stars", 1, "Show repos if they have this quantity of stars.")

	NpmOwners = flag.String("npm-owner", "", "List of NPM owners whose packages to query (comma separated).")
	NpmWeekly = flag.Int("npm-weekly", 1, "Show repos if they have this quantity of weekly downloads.")
)

func init() {
	flag.Parse()
}
