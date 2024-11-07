package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Shackelford-Arden/hctx/build"
	"github.com/Shackelford-Arden/hctx/internal/github"
)

func ShowVersion(ctx *cli.Context) error {

	fmt.Printf("Current version: %s\n", build.Version)

	if ctx.Bool("check-latest") {

		gh := github.NewClient("", "")
		gt := os.Getenv("GITHUB_TOKEN")

		if gt != "" {
			gh.SetToken(gt)
		}

		latestVersion, lvError := gh.GetLatestRelease()
		if lvError != nil {
			return fmt.Errorf("failed to get latest version: %s", lvError.Error())
		}

		fmt.Printf("Latest version: %s\n", latestVersion.TagName)

	}

	return nil

}
