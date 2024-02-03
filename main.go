package main

import (
	"fmt"
	"github.com/Shackelford-Arden/hctx/cmd"
	"log"
	"os"
)

func main() {

	app, appErr := cmd.App()
	if appErr != nil {
		fmt.Printf("error running hctx: %s", appErr)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
