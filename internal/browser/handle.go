package browser

import (
	"context"
	"log"
	"path"

	"github.com/wmalik/ogit/internal/db"
	"github.com/wmalik/ogit/internal/gitconfig"
	"github.com/wmalik/ogit/internal/gitutils"
	"github.com/wmalik/ogit/internal/sync"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleCommandDefault(useCache bool) error {
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

	if !useCache {
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
