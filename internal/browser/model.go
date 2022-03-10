package browser

import (
	"fmt"
	"sort"
	"time"

	"github.com/wmalik/ogit/internal/db"
	"github.com/wmalik/ogit/internal/gitutils"
	"github.com/wmalik/ogit/service"

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
	// the storage path of the selected item
	selectedItemStoragePath string
	// whether a shell should be spawned after the TUI exits
	spawnShell bool

	gu *gitutils.GitUtils
	rs *service.RepositoryService
}

func NewModelWithItems(repos []db.Repository, storagePath string, gu *gitutils.GitUtils) *model {

	listItems := sortItemsCloned(toItems(repos, storagePath))
	m := list.NewModel(listItems, delegateItemUpdate(storagePath), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[ogit] [%s]", storagePath)
	m.Styles.Title = titleBarStyle
	m.AdditionalShortHelpKeys = availableKeyBindingsCB
	m.SetShowStatusBar(false)

	return &model{
		list:            m,
		storagePath:     storagePath,
		bottomStatusBar: "-",
		gu:              gu,
	}
}

func toItems(repos []db.Repository, storagePath string) []list.Item {
	items := make([]list.Item, len(repos))

	for i := range repos {
		repoItem := newRepoItem(&repos[i], storagePath)
		if repoItem.Cloned() {
			repoItem.SetTitle(brightStyle.Render(repoItem.Repository.Title))
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
		return items[i].(repoItem).Repository.Title < items[j].(repoItem).Repository.Title
	})
	return items
}
