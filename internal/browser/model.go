package browser

import (
	"fmt"
	"ogit/service"
	"strings"
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

func NewModel(orgs []string, cloneDirPath string, repoService *service.RepositoryService) model {
	// Start with an empty list of items
	m := list.NewModel([]list.Item{}, delegateItemUpdate(cloneDirPath, orgs, repoService), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[Repositories] [%s] [%s]", strings.Join(orgs, " "), cloneDirPath)
	m.AdditionalShortHelpKeys = availableKeyBindingsCB

	return model{
		list:            m,
		orgs:            orgs,
		cloneDirPath:    cloneDirPath,
		rs:              repoService,
		bottomStatusBar: "-",
	}
}
