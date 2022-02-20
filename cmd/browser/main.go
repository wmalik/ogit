package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"ogit/internal/browser"
	"ogit/internal/db"
	"ogit/internal/gitconfig"
	"ogit/internal/gitutils"
	"ogit/internal/sync"

	tea "github.com/charmbracelet/bubbletea"
)

const usageTpl = `
Usage: %s [OPTION]
Organize git repositories
Sync repositories on startup unless -nosync is specified
`

func usage() func() {
	return func() {
		rendered := fmt.Sprintf(usageTpl, os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), rendered)
		flag.PrintDefaults()
	}
}

func main() {
	noSync := flag.Bool("nosync", false, "Disable syncing of repositories metadata at startup")
	clear := flag.Bool("clear", false, "Clear all local repository metadata")

	flag.Usage = usage()

	flag.Parse()

	ctx := context.Background()

	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	localDB, err := db.NewDB(path.Join(gitConf.CloneDirPath(), "ogit.db"))
	if err != nil {
		log.Fatalln(err)
	}

	if err := localDB.Init(); err != nil {
		log.Fatalln(err)
	}

	if *clear {
		if err := localDB.DeleteAllRepositories(ctx); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	if *noSync == false {
		if err := sync.Sync(ctx, gitConf); err != nil {
			log.Fatalln(err)
		}
	}

	repos, err := localDB.SelectAllRepositories(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	gu, err := gitutils.NewGitUtils(gitConf.UseSSHAgent(), gitConf.PrivKeyPath())
	if err != nil {
		log.Fatalln(err)
	}

	f, err := tea.LogToFile("/tmp/ogit.log", "debug")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := tea.NewProgram(
		browser.NewModelWithItems(repos, gitConf.CloneDirPath(), gu),
	).Start(); err != nil {
		log.Fatalln(err)
	}
}
