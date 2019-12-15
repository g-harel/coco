package github

import (
	"fmt"
	"net/http"

	"github.com/g-harel/coco/internal/httpc"
)

func fetchGeneric(url, token string, body interface{}) (*http.Header, error) {
	return httpc.Get(
		url,
		http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}},
		body,
	)
}

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
