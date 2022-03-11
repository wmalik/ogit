package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wmalik/ogit/internal/browser"
	"github.com/wmalik/ogit/internal/bulkclone"
	"github.com/wmalik/ogit/internal/clear"

	"github.com/urfave/cli/v2"
)

// ldflags populated by goreleaser
var (
	version = "development"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	app := &cli.App{
		Usage:       "Organize git repositories",
		Version:     fmt.Sprintf("%s %s %s", version, commit, date),
		HideVersion: false,
		Action: func(c *cli.Context) error {
			if err := browser.HandleCommandDefault(); err != nil {
				log.Fatalln(err)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Usage:   "Fetch repository metadata",
				Action: func(c *cli.Context) error {
					if err := browser.HandleCommandFetch(); err != nil {
						log.Fatalln(err)
					}
					return nil
				},
			},
			{
				Name:    "clone",
				Aliases: []string{"c"},
				Usage:   "Clone repositories in bulk",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "org",
						Usage:    "Organization name",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "filter",
						Usage: "filter repositories by name",
					},
				},
				Action: func(c *cli.Context) error {
					if err := bulkclone.HandleCommandClone(c.String("org"), c.String("filter")); err != nil {
						log.Fatalln(err)
					}
					return nil
				},
			},
			{
				Name:  "clear",
				Usage: "Clear all local repository metadata (not the repository contents)",
				Action: func(c *cli.Context) error {
					if err := clear.HandleCommandDefault(c.Context); err != nil {
						log.Fatalln(err)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
