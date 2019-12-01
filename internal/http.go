package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v", res.StatusCode)
	}

	data, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(data), res.Header)

	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return nil, err
	}

	return &res.Header, nil
}
