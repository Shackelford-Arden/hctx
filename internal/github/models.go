package github

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Asset struct {
	Url                string    `json:"url"`
	Id                 int       `json:"id"`
	NodeId             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

func (r *Release) Version() string {
	return r.TagName
}

func (r *Release) TarballUrl() string {

	fmt.Println(fmt.Sprintf("Looking for %s", releaseName()))
	for _, asset := range r.Assets {
		fmt.Println(asset.Name)
		if strings.HasSuffix(asset.Name, releaseName()) {
			return asset.BrowserDownloadUrl
		}
	}

	return ""

}

func (r *Release) ChecksumUrl() string {

	for _, asset := range r.Assets {
		if strings.HasSuffix(asset.Name, "checksums.txt") {
			return asset.BrowserDownloadUrl
		}
	}

	return ""

}

func releaseName() string {

	// hctx's goreleaser config uses names like Darwin/Linux/Windows, so
	// we need to uppercase the first letter.
	osName := runtime.GOOS
	osName = strings.ToUpper(osName[:1]) + osName[1:]

	osArch := runtime.GOARCH

	return fmt.Sprintf("hctx_%s_%s.tar.gz", osName, osArch)

}
