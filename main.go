package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shackelford-Arden/hctx/cmd"
)

func main() {

	app, appErr := cmd.App()
	if appErr != nil {
		fmt.Printf("error running hctx: %s", appErr)
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
