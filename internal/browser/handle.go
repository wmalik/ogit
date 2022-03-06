package browser

import (
	"context"
	"log"
	"ogit/internal/db"
	"ogit/internal/gitconfig"
	"ogit/internal/gitutils"
	"ogit/internal/sync"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleCommandDefault(noSync, clear bool) error {
	ctx := context.Background()
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	localDB, err := db.NewDB(path.Join(gitConf.StoragePath(), "ogit.db"))
	if err != nil {
		log.Fatalln(err)
	}

	if err := localDB.Init(); err != nil {
		log.Fatalln(err)
	}

	if clear {
		if err := localDB.DeleteAllRepositories(ctx); err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)
	}

	if noSync == false {
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
		NewModelWithItems(repos, gitConf.StoragePath(), gu),
	).Start(); err != nil {
		log.Fatalln(err)
	}
	return nil
}
