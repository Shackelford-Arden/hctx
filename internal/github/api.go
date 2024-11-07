package github

import (
	"net/http"
)

type LatestVersion struct {
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

func (c *Client) GetLatestRelease() (*Release, error) {
	path := "/repos/Shackelford-Arden/hctx/releases/latest"

	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var release Release
	if err := c.doRequest(req, &release); err != nil {
		return nil, err
	}

	return &release, nil
}
