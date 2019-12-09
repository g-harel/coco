package flags

import (
	"flag"
)

var (
	GithubToken  = flag.String("github-token", "", "GitHub API token.")
	GithubOwners = flag.String("github-owner", "", "List of GitHub owners whose repos to query (comma separated).")
	GithubViews  = flag.Int("github-views", 0, "Show repos if they have this quantity of views.")
	GithubToday  = flag.Int("github-today", 0, "Show repos if they have this quantity of views today.")

	NpmOwners = flag.String("npm-owner", "", "List of NPM owners whose packages to query (comma separated).")
	NpmWeekly = flag.Int("npm-weekly", 0, "Show repos if they have this quantity of weekly downloads.")

	LogErrors = flag.Bool("log-error", true, "Log errors.")
	LogInfo   = flag.Bool("log-info", false, "Log miscellaneous info.")
)

func init() {
	flag.Parse()
}
