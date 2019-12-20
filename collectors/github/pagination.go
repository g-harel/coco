package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GeneratePaginatedURLs is a helper to navigate the GitHub
// API's pagination scheme. It uses the headers from an
// initial response to generate the remaining URLs to fetch.
func generatePaginatedURLs(h *http.Header) ([]string, error) {
	links := strings.Split(h.Get("Link"), ",")
	// Find URL marked as last.
	rawLastURL := ""
	for i := 0; i < len(links); i++ {
		if strings.HasSuffix(links[i], `>; rel="last"`) {
			rawURL := strings.TrimSpace(links[i])
			rawURL = strings.TrimPrefix(rawURL, "<")
			rawLastURL = strings.TrimSuffix(rawURL, `>; rel="last"`)
		}
	}
	if rawLastURL == "" {
		// Pagination header is missing when the response
		// contains all data.
		return []string{rawLastURL}, nil
	}
	lastURL, err := url.Parse(rawLastURL)
	if err != nil {
		return nil, fmt.Errorf("parse pagination url: %v", err)
	}
	lastPageIndex, err := strconv.Atoi(lastURL.Query().Get("page"))
	if err != nil {
		return nil, fmt.Errorf("parse last pagination index: %v", err)
	}
	urls := []string{}
	for i := 1; i <= lastPageIndex; i++ {
		nthPageURL, _ := url.Parse(lastURL.String())
		// Replace query's page value.
		query := nthPageURL.Query()
		query.Del("page")
		query.Add("page", strconv.Itoa(i))
		nthPageURL.RawQuery = query.Encode()
		urls = append(urls, nthPageURL.String())
	}
	return urls, nil
}
