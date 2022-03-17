package repocommands

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/wmalik/ogit/internal/db"
	"github.com/wmalik/ogit/internal/gitconfig"
	"github.com/wmalik/ogit/internal/utils"
)

type Command int32

// Various command types supported by HandleURLCommands.
const (
	Pulls Command = iota
	Web
	Org
	Issues
	CI
	Releases
	Settings
)

// HandleURLCommands opens the relevant URL in the web browser.
func HandleURLCommands(ctx context.Context, command Command) error {
	repo, err := findRepositoryCWD(ctx)
	if err != nil {
		if err.Error() == "record not found" {
			fmt.Println("error: repository not managed by ogit")
			return nil
		}
		return err
	}

	var url string

	switch command {
	case Pulls:
		url = repo.BrowserPullRequestsURL
	case Web:
		url = repo.BrowserHomepageURL
	case Org:
		url = repo.OrgURL
	case Issues:
		url = repo.IssuesURL
	case CI:
		url = repo.CIURL
	case Releases:
		url = repo.ReleasesURL
	case Settings:
		url = repo.SettingsURL
	}

	if err := utils.OpenURL(url); err != nil {
		return err
	}

	fmt.Println(url)

	return nil
}

func findRepositoryCWD(ctx context.Context) (*db.Repository, error) {
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		return nil, err
	}

	database, err := db.NewDB(path.Join(gitConf.StoragePath(), "ogit.db"))
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := database.FindRepository(ctx,
		providerFromPath(cwd), orgFromPath(cwd), nameFromPath(cwd))
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func providerFromPath(path string) string {
	return filepath.Base(filepath.Dir(filepath.Dir(path)))
}

func orgFromPath(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func nameFromPath(path string) string {
	return filepath.Base(path)
}
