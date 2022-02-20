package browser

import (
	"fmt"
	"ogit/internal/gitutils"
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

func NewModel(orgs []string, gitlabGroups []string, cloneDirPath string, repoService *service.RepositoryService, gu *gitutils.GitUtils) model {
	// Start with an empty list of items
	m := list.NewModel([]list.Item{}, delegateItemUpdate(cloneDirPath, orgs, gitlabGroups, repoService, gu), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[Repositories] [%s] [%s] [%s]", strings.Join(orgs, " "), strings.Join(gitlabGroups, " "), cloneDirPath)
	m.AdditionalShortHelpKeys = availableKeyBindingsCB

	return model{
		list:            m,
		orgs:            orgs,
		cloneDirPath:    cloneDirPath,
		rs:              repoService,
		bottomStatusBar: "-",
	}
}
