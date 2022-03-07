package browser

import (
	"fmt"
	"ogit/internal/db"
	"ogit/internal/gitutils"
	"ogit/service"
	"path"
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
	storagePath string
	// A status bar to show useful information e.g. Github API usage
	bottomStatusBar string

	rs *service.RepositoryService
}

func NewModelWithItems(repos []db.Repository, storagePath string, gu *gitutils.GitUtils) model {

	listItems := sortItemsCloned(toItems(repos, storagePath))
	m := list.NewModel(listItems, delegateItemUpdate(storagePath, gu), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[ogit] [%s]", storagePath)
	m.Styles.Title = titleBarStyle
	m.AdditionalShortHelpKeys = availableKeyBindingsCB
	m.SetShowStatusBar(false)

	return model{
		list:            m,
		storagePath:    storagePath,
		bottomStatusBar: "-",
	}
}

func toItems(repos []db.Repository, storagePath string) []list.Item {
	items := make([]list.Item, len(repos))

	for i := range repos {
		repoItem := repoItem{
			title:                  repos[i].Title,
			owner:                  repos[i].Owner,
			name:                   repos[i].Name,
			description:            repos[i].Description,
			browserHomepageURL:     repos[i].BrowserHomepageURL,
			browserPullRequestsURL: repos[i].BrowserPullRequestsURL,
			httpsCloneURL:          repos[i].HTTPSCloneURL,
			sshCloneURL:            repos[i].SSHCloneURL,
			storagePath:            path.Join(storagePath, repos[i].Owner, repos[i].Name),
		}

		if repoItem.Cloned() {
			repoItem.title = brightStyle.Render(repoItem.title)
		}
		items[i] = repoItem
	}

	return items
}

func sortItemsCloned(items []list.Item) []list.Item {
	// sort items by whether they have been cloned
	sort.Slice(items, func(i, j int) bool {
		return items[i].(repoItem).Cloned()
	})

	// sort items in lexical order
	sort.Slice(items, func(i, j int) bool {
		return items[i].(repoItem).Title() < items[j].(repoItem).Title()
	})
	return items
}
