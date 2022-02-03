package apiserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	httpClient
	apiURL string
	apiKey string
}

func newClient(apiURL, apiKey string) *Client {
	return &Client{
		http.DefaultClient,
		apiURL,
		apiKey,
	}
}

func (c *Client) get(params map[string]string, targetModel interface{}) (statusCode int, err error) {
	query := c.buildQuery(params)
	url, err := url.Parse(c.apiURL)
	if err != nil {
		return 0, err
	}

	url.RawQuery = query
	resp, err := c.Get(url.String())
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	switch resp.StatusCode {
	case 200:
		if targetModel != nil {
			err = json.Unmarshal(body, targetModel)
			if err != nil {
				return 0, err
			}
		}
		return resp.StatusCode, nil

	default:
		return resp.StatusCode, fmt.Errorf("%s", body)
	}
}

func (c *Client) buildQuery(params map[string]string) string {
	v := url.Values{}
	for key, value := range params {
		v.Set(key, value)
	}

	v.Set("apikey", c.apiKey)

	return v.Encode()
}
