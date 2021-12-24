package browser

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
)

type model struct {
	list  list.Model
	fetch bool
}

type repoListItem struct {
	title       string
	description string
	browserURL  string
	cloneURL    string
}

type listKeyMap struct {
	toggleHelpMenu key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func NewModel() model {
	// Start with an empty list of items
	reposList := list.NewModel([]list.Item{}, repoListItemDelegate(), 0, 0)
	reposList.Title = "Repositories"
	reposList.Styles.Title = titleStyle
	reposList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("H", "R"),
				key.WithHelp("H", "toggle help"),
				key.WithHelp("R", "Refresh list"),
			),
		}
	}

	return model{
		list:  reposList,
		fetch: true,
	}
}

func (i repoListItem) Title() string       { return i.title }
func (i repoListItem) Description() string { return i.description }
func (i repoListItem) FilterValue() string { return i.title + i.description }
func (i repoListItem) BrowserURL() string  { return i.browserURL }
func (i repoListItem) CloneURL() string    { return i.cloneURL }

func repoListItemDelegate() list.DefaultDelegate {
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
	d.UpdateFunc = delegateUpdateFunc(keyBinding)

	return d
}
