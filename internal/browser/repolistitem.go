package browser

import (
	"context"
	"io"
	"log"
	"ogit/internal/gitutils"
)

type repoItem struct {
	title                  string
	owner                  string
	name                   string
	description            string
	browserHomepageURL     string
	browserPullRequestsURL string
	httpsCloneURL          string
	sshCloneURL            string
	storagePath            string
}

func (i repoItem) Title() string                  { return i.title }
func (i repoItem) Owner() string                  { return i.owner }
func (i repoItem) Name() string                   { return i.name }
func (i repoItem) Description() string            { return i.description }
func (i repoItem) FilterValue() string            { return i.title + i.description }
func (i repoItem) BrowserHomepageURL() string     { return i.browserHomepageURL }
func (i repoItem) BrowserPullRequestsURL() string { return i.browserPullRequestsURL }
func (i repoItem) HTTPSCloneURL() string          { return i.httpsCloneURL }
func (i repoItem) SSHCloneURL() string            { return i.sshCloneURL }
func (i repoItem) StoragePath() string            { return i.storagePath }
func (i repoItem) Cloned() bool {
	return gitutils.Cloned(i.storagePath)
}

type cloneService interface {
	CloneToDisk(ctx context.Context, httpsURL string, sshURL string, path string, progress io.Writer) (string, error)
}

func (i repoItem) Clone(ctx context.Context, cloner cloneService) (string, error) {
	return cloner.CloneToDisk(ctx, i.httpsCloneURL, i.sshCloneURL, i.storagePath, log.Default().Writer())
}

func (i repoItem) LastCommitInfo() (string, error) {
	if i.Cloned() {
		repo, err := gitutils.ReadRepository(i.storagePath)
		if err != nil {
			return "", err
		}

		return repo.LastCommit(), nil
	}

	return "", nil
}
