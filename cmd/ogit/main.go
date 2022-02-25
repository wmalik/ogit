package main

import (
	"log"
	"os"

	"ogit/internal/browser"
	"ogit/internal/bulkclone"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "Organize git repositories",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "nosync",
				Usage: "Disable syncing of repositories metadata at startup",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "clear",
				Usage: "Clear all local repository metadata",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			if err := browser.HandleCommandDefault(c.Bool("nosync"), c.Bool("clear")); err != nil {
				log.Fatalln(err)
			}
			return nil
		},
		Commands: []*cli.Command{
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
