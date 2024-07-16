package cmd

import (
	"fmt"

	"github.com/Shackelford-Arden/hctx/build"
	"github.com/urfave/cli/v2"
)

func ShowVersion(ctx *cli.Context) error {
	fmt.Println(fmt.Sprintf("%s %s - Commit %s & built with %s on %s", ctx.App.Name, ctx.App.Version, build.Commit, build.BuiltWith, build.Date))
	return nil
}
