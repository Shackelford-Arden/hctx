package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	token string

	addr       string
	httpClient *http.Client
}

const baseURL = "https://api.github.com"

func NewClient(token string, addr string) *Client {
	if addr == "" {
		addr = "https://api.github.com"
	}

	return &Client{
		addr:       addr,
		httpClient: &http.Client{},
		token:      token,
	}
}

// newRequest creates a new http.Request object
func (c *Client) newRequest(method, path string, queryParams url.Values) (*http.Request, error) {
	reqUrl := fmt.Sprintf("%s%s", baseURL, path)
	if queryParams != nil {
		reqUrl += "?" + queryParams.Encode()
	}

	req, err := http.NewRequest(method, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

func (c *Client) SetToken(token string) {
	c.token = token
}
