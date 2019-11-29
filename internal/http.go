package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func HTTPGet(rawUrl string, headers http.Header, body interface{}) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	res, err := DefaultLoggingClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"x-spiferack": []string{"1"},
		},
	})
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %v", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return err
	}

	return nil
}
