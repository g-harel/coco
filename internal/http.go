package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func HTTPGet(rawUrl string, headers http.Header, body interface{}) (*http.Header, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	res, err := DefaultLoggingClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: headers,
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return nil, err
	}

	return &res.Header, nil
}
