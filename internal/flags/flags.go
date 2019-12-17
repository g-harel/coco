package flags

import (
	"flag"
)

// Global list of all flags understood by the command.
var (
	LogInfo   = flag.Bool("log-info", false, "Log miscellaneous info.")
	LogErrors = flag.Bool("log-error", true, "Log errors.")

	RateLimit = flag.Int("rate-limit", 999, "Maximum number of requests per second by all collectors.")

	GithubOwners = multiFlag{}
	GithubToken  = flag.String("github-token", "", "GitHub API token.")
	GithubViews  = flag.Int("github-views", 1, "Show repos if they have this quantity of views.")
	GithubToday  = flag.Int("github-today", 1, "Show repos if they have this quantity of views today.")
	GithubStars  = flag.Int("github-stars", 1, "Show repos if they have this quantity of stars.")

	NpmOwners = multiFlag{}
	NpmWeekly = flag.Int("npm-weekly", 1, "Show repos if they have this quantity of weekly downloads.")
)

func init() {
	flag.Var(&GithubOwners, "github-owner", "GitHub owner to query.")
	flag.Var(&NpmOwners, "npm-owner", "NPM owners whose packages to query.")
	flag.Parse()
}
