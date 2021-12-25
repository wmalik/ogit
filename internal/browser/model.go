package browser

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// The state of browser
type model struct {
	// the list of repositories
	list list.Model
	// whether the list should be fetched from remote
	fetch bool
	// the list of github organisations
	orgs []string
	// the path on disk where repositories should be cloned
	cloneDirPath string
}

type listKeyMap struct {
	toggleHelpMenu key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func NewModel(orgs []string, cloneDirPath string) model {
	// Start with an empty list of items
	m := list.NewModel([]list.Item{}, repoListItemDelegate(cloneDirPath), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[Repositories] [%s] [%s]", strings.Join(orgs, " "), cloneDirPath)
	m.Styles.Title = titleStyle
	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("H", "R"),
				key.WithHelp("H", "toggle help"),
				key.WithHelp("R", "Refresh list"),
			),
		}
	}

	return model{
		list:         m,
		fetch:        true,
		orgs:         orgs,
		cloneDirPath: cloneDirPath,
	}
}

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

func repoListItemDelegate(cloneDirPath string) list.DefaultDelegate {
	keyBinding := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
		key.WithKeys("o"),
		key.WithHelp("o", "open in firefox"),
		key.WithKeys("c"),
		key.WithHelp("c", "clone repository"),
	)

	d := list.NewDefaultDelegate()
	d.ShortHelpFunc = func() []key.Binding {
		return []key.Binding{keyBinding}
	}
	d.UpdateFunc = delegateUpdateFunc(keyBinding, cloneDirPath)

	return d
}
