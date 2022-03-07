package browser

import (
    "context"
    "io"
    "log"
	"ogit/internal/gitutils"
)

type repoListItem struct {
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

func (i repoListItem) Title() string                  { return i.title }
func (i repoListItem) Owner() string                  { return i.owner }
func (i repoListItem) Name() string                   { return i.name }
func (i repoListItem) Description() string            { return i.description }
func (i repoListItem) FilterValue() string            { return i.title + i.description }
func (i repoListItem) BrowserHomepageURL() string     { return i.browserHomepageURL }
func (i repoListItem) BrowserPullRequestsURL() string { return i.browserPullRequestsURL }
func (i repoListItem) HTTPSCloneURL() string          { return i.httpsCloneURL }
func (i repoListItem) SSHCloneURL() string            { return i.sshCloneURL }
func (i repoListItem) StoragePath() string            { return i.storagePath }
func (i repoListItem) Cloned() bool {
	return gitutils.Cloned(i.storagePath)
}

type cloneService interface {
  CloneToDisk(ctx context.Context, httpsURL string, sshURL string, path string, progress io.Writer) (string, error)
}

func (i repoListItem) Clone(ctx context.Context, cloner cloneService) (string, error) {
  return cloner.CloneToDisk(ctx, i.httpsCloneURL, i.sshCloneURL, i.storagePath, log.Default().Writer())
}

func (i repoListItem) LastCommitInfo() (string, error) {
	if i.Cloned() {
		repo, err := gitutils.ReadRepository(i.storagePath)
		if err != nil {
			return "", err
		}

		return repo.LastCommit(), nil
	}

	return "", nil
}
