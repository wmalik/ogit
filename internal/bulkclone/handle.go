package bulkclone

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/wmalik/ogit/internal/db"
	"github.com/wmalik/ogit/internal/gitconfig"
	"github.com/wmalik/ogit/internal/gitutils"

	"github.com/charmbracelet/lipgloss"
)

func HandleCommandClone(org, filter string) error {
	ctx := context.Background()
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		return err
	}

	localDB, err := db.NewDB(path.Join(gitConf.StoragePath(), "ogit.db"))
	if err != nil {
		return err
	}

	if err := localDB.Init(); err != nil {
		return err
	}

	repos, err := localDB.SelectRepositories(ctx, org, filter)
	if err != nil {
		return err
	}

	gu, err := gitutils.NewGitUtils(gitConf.UseSSHAgent(), gitConf.PrivKeyPath())
	if err != nil {
		return err
	}

	cloneFailed := false
	for _, repo := range repos {
		clonePath := path.Join(gitConf.StoragePath(), repo.Provider, repo.Owner, repo.Name)

		if gitutils.Cloned(clonePath) {
			printMsgDimmed(fmt.Sprintf("[already cloned] %s/%s", repo.Owner, repo.Name))
			continue
		}
		printMsg(fmt.Sprintf("Cloning %s/%s", repo.Owner, repo.Name))
		_, err := gu.CloneToDisk(context.Background(),
			repo.HTTPSCloneURL,
			repo.SSHCloneURL,
			clonePath,
			ioutil.Discard,
		)
		if err != nil {
			printMsg(fmt.Sprintf("unable to clone %s/%s %s\n", repo.Owner, repo.Name, err))
			cloneFailed = true
		}
	}

	if cloneFailed {
		return fmt.Errorf("failed to clone one or more repos")
	}
	return nil
}

func printMsg(message string) {
	fmt.Printf("* %s\n", message)
}

func printMsgDimmed(message string) {
	printMsg(lipgloss.NewStyle().Faint(true).Render(message))
}
