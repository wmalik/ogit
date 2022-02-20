package browser

import (
	"fmt"
	"ogit/internal/db"
	"ogit/internal/gitutils"
	"ogit/service"
	"time"

	"github.com/charmbracelet/bubbles/list"
)

// The state of browser
type model struct {
	// the list of repositories
	list list.Model
	// list of organisations or users (currently only public users or organisations)
	orgs []string
	// the path on disk where repositories should be cloned
	cloneDirPath string
	// A status bar to show useful information e.g. Github API usage
	bottomStatusBar string

	rs *service.RepositoryService
}

func NewModelWithItems(repos []db.Repository, cloneDirPath string, gu *gitutils.GitUtils) model {

	items := make([]list.Item, len(repos))

	for i := range repos {
		repoItem := repoListItem{
			title:                  repos[i].Title,
			owner:                  repos[i].Owner,
			name:                   repos[i].Name,
			description:            repos[i].Description,
			browserHomepageURL:     repos[i].BrowserHomepageURL,
			browserPullRequestsURL: repos[i].BrowserPullRequestsURL,
			httpsCloneURL:          repos[i].HTTPSCloneURL,
			sshCloneURL:            repos[i].SSHCloneURL,
		}

		if repoItem.Cloned(cloneDirPath) {
			repoItem.title = statusMessageStyle(repoItem.Title())
			repoItem.description = statusMessageStyle(repoItem.Description())
		}
		items[i] = repoItem
	}
	m := list.NewModel(items, delegateItemUpdate(cloneDirPath, gu), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[Repositories] [%s]", cloneDirPath)
	m.AdditionalShortHelpKeys = availableKeyBindingsCB

	return model{
		list:            m,
		cloneDirPath:    cloneDirPath,
		bottomStatusBar: "-",
	}
}
