package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s - version %s\n", c.App.Name, version)
		fmt.Printf("  commit: \t%s\n", commit)
		fmt.Printf("  build date: \t%s\n", date)
		fmt.Printf("  build user: \t%s\n", builtBy)
		fmt.Printf("  go version: \t%s\n", runtime.Version())
	}

	app := &cli.App{
		Name:    "golang-repo-template",
		Usage:   "golang-repo-template usage",
		Version: version,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// These values are private which ensures they can only be set with the build flags.
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)
