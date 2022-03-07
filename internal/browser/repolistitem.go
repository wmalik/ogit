package browser

import (
	"context"
	"io"
	"log"
	"path"

	"github.com/wmalik/ogit/internal/db"
	"github.com/wmalik/ogit/internal/gitutils"
)

type repoItem struct {
	*db.Repository
	repoStoragePath string
}

func newRepoItem(repo *db.Repository, storageBasePath string) repoItem {
	return repoItem{
		Repository:      repo,
		repoStoragePath: path.Join(storageBasePath, repo.Provider, repo.Owner, repo.Name),
	}
}

func (i repoItem) Title() string       { return i.Repository.Title }
func (i repoItem) Description() string { return i.Repository.Description }
func (i repoItem) FilterValue() string { return i.Repository.Title + i.Repository.Description }
func (i repoItem) StoragePath() string {
	return i.repoStoragePath
}
func (i repoItem) Cloned() bool {
	return gitutils.Cloned(i.StoragePath())
}

func (i *repoItem) SetTitle(title string) { i.Repository.Title = title }

type cloneService interface {
	CloneToDisk(ctx context.Context, httpsURL string, sshURL string, path string, progress io.Writer) (string, error)
}

func (i repoItem) Clone(ctx context.Context, cloner cloneService) (string, error) {
	return cloner.CloneToDisk(ctx, i.Repository.HTTPSCloneURL, i.Repository.SSHCloneURL, i.StoragePath(), log.Default().Writer())
}
