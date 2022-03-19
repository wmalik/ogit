package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wmalik/ogit/internal/browser"
	"github.com/wmalik/ogit/internal/bulkclone"
	"github.com/wmalik/ogit/internal/clear"
	"github.com/wmalik/ogit/internal/repocommands"

	"github.com/urfave/cli/v2"
)

// ldflags populated by goreleaser.
var (
	version = ""
	commit  = ""
	date    = ""
)

func main() {
	app := &cli.App{
		Usage:       "Organize git repositories",
		Version:     fmt.Sprintf("%s %s %s", version, commit, date),
		HideVersion: false,
		Action: func(c *cli.Context) error {
			if err := browser.HandleCommandDefault(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "fetch",
				Usage: "Fetch all repository metadata from provider APIs (e.g. GitHub/GitLab)",
				Action: func(c *cli.Context) error {
					if err := browser.HandleCommandFetch(); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "clone",
				Usage: "Clone repositories of an organization",
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
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "clear",
				Usage: "Clear all local repository metadata (not the repository contents)",
				Action: func(c *cli.Context) error {
					if err := clear.HandleCommandDefault(c.Context); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "pulls",
				Aliases: []string{"prs", "mrs"},
				Usage:   "Open repository pull requests in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Pulls); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "web",
				Aliases: []string{"home"},
				Usage:   "Open repository home page in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Web); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "org",
				Usage: "Open repository org in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Org); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "issues",
				Usage: "Open repository issues in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Issues); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:    "ci",
				Aliases: []string{"actions"},
				Usage:   "Open repository CI/actions in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.CI); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "releases",
				Usage: "Open repository releases in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Releases); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
			{
				Name:  "settings",
				Usage: "Open repository settings in web browser",
				Action: func(c *cli.Context) error {
					if err := repocommands.HandleURLCommands(c.Context, repocommands.Settings); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
