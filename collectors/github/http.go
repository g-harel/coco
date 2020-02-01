package github

import (
	"fmt"
	"net/http"

	"github.com/g-harel/coco/internal/httpc"
)

// FetchGeneric is a shared wrapper that adds the GitHub
// Authorization header to requests.
func fetchGeneric(url, token string, body interface{}) (*http.Header, error) {
	return httpc.Get(
		url,
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}},
		body,
	)
}

// FetchFirstRepo fetches repo lists from the GitHub API.
// It can only be used for the first request because the
// pagination scheme is sent back as headers and is not
// compatible with simply increasing the page number on this
// URL.
func fetchFirstRepos(token, owner string) (reposResponse, *http.Header, error) {
	res := reposResponse{}
	h, err := fetchGeneric(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=1", owner),
		token,
		&res,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch owner initial page %v: %v", owner, err)
	}
	return res, h, nil

}

// FetchViews fetches repo views data from the GitHub API.
func fetchViews(token, owner, name string) (*viewsResponse, error) {
	res := &viewsResponse{}
	_, err := fetchGeneric(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		token,
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch repo views %v: %v", name, err)
	}
	return res, nil
}

func fetchUser(token, user string) (*userResponse, error) {
	res := &userResponse{}
	_, err := fetchGeneric(
		fmt.Sprintf("https://api.github.com/users/%v", user),
		token,
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch user info %v: %v", user, err)
	}
	return res, nil
}
