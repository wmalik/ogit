package browser

import (
	"fmt"
	"ogit/internal/db"
	"ogit/internal/gitutils"
	"ogit/service"
	"sort"
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

	listItems := sortItemsCloned(toItems(repos, cloneDirPath), cloneDirPath)
	m := list.NewModel(listItems, delegateItemUpdate(cloneDirPath, gu), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[ogit] [%s]", cloneDirPath)
	m.Styles.Title = titleBarStyle
	m.AdditionalShortHelpKeys = availableKeyBindingsCB
	m.SetShowStatusBar(false)

	return model{
		list:            m,
		cloneDirPath:    cloneDirPath,
		bottomStatusBar: "-",
	}
}

func toItems(repos []db.Repository, cloneDirPath string) []list.Item {
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
			repoItem.title = brightStyle.Render(repoItem.title)
		}
		items[i] = repoItem
	}

	return items
}

func sortItemsCloned(items []list.Item, cloneDirPath string) []list.Item {
	// sort items by whether they have been cloned
	sort.Slice(items, func(i, j int) bool {
		return items[i].(repoListItem).Cloned(cloneDirPath)
	})

	// sort items in lexical order
	sort.Slice(items, func(i, j int) bool {
		return items[i].(repoListItem).Title() < items[j].(repoListItem).Title()
	})
	return items
}
