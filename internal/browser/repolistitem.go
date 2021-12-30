package browser

import (
	"ogit/internal/gitutils"
	"path"
)

type repoListItem struct {
	title       string
	owner       string
	name        string
	description string
	browserURL  string
	cloneURL    string
}

func (i repoListItem) Title() string       { return i.title }
func (i repoListItem) Owner() string       { return i.owner }
func (i repoListItem) Name() string        { return i.name }
func (i repoListItem) Description() string { return i.description }
func (i repoListItem) FilterValue() string { return i.title + i.description }
func (i repoListItem) BrowserURL() string  { return i.browserURL }
func (i repoListItem) CloneURL() string    { return i.cloneURL }
func (i repoListItem) Cloned(cloneDirPath string) bool {
	return gitutils.Cloned(path.Join(cloneDirPath, i.owner, i.name))
}

func (i repoListItem) LastCommitInfo(cloneDirPath string) (string, error) {
	if i.Cloned(cloneDirPath) {
		repo, err := gitutils.ReadRepository(path.Join(cloneDirPath, i.owner, i.name))
		if err != nil {
			return "", err
		}

		return repo.LastCommit(), nil
	}

	return "", nil
}
