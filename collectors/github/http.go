package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/g-harel/coco/internal/httpc"
)

func parsePaginationHeaderLastURL(h *http.Header) (*url.URL, error) {
	header := strings.Split(h.Get("Link"), ",")
	for i := 0; i < len(header); i++ {
		if strings.HasSuffix(header[i], `>; rel="last"`) {
			rawURL := strings.TrimSpace(header[i])
			rawURL = strings.TrimPrefix(rawURL, "<")
			rawURL = strings.TrimSuffix(rawURL, `>; rel="last"`)
			url, err := url.Parse(rawURL)
			if err != nil {
				return nil, fmt.Errorf("parse pagination url: %v", err)
			}
			return url, nil

		}
	}
	return nil, fmt.Errorf("pagination header value not found")
}

func fetchInitialOwnerRepo(token, owner string) (repoListResponse, *url.URL, error) {
	res := repoListResponse{}
	h, err := httpc.Get(
		fmt.Sprintf("https://api.github.com/users/%v/repos?page=1", owner),
		tokenHeader(token),
		&res,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("fetch owner initial page %v: %v", owner, err)
	}
	lastPageURL, err := parsePaginationHeaderLastURL(h)
	if err != nil {
		return res, nil, nil
	}
	return res, lastPageURL, nil

}

func fetchRepoViews(token, owner, name string) (*repoViewsResponse, error) {
	res := &repoViewsResponse{}
	_, err := httpc.Get(
		fmt.Sprintf("https://api.github.com/repos/%v/%v/traffic/views", owner, name),
		tokenHeader(token),
		res,
	)
	if err != nil {
		return nil, fmt.Errorf("fetch repo views %v: %v", name, err)
	}
	return res, nil
}

func tokenHeader(token string) http.Header {
	return http.Header{"Authorization": []string{fmt.Sprintf("token %v", token)}}
}
