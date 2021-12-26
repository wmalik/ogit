package browser

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
)

// The state of browser
type model struct {
	// the list of repositories
	list list.Model
	// whether the list should be fetched from remote
	// TODO move this to a tea.Cmd
	fetch bool
	// the list of github organisations
	orgs []string
	// the path on disk where repositories should be cloned
	cloneDirPath string
}

func NewModel(orgs []string, cloneDirPath string) model {
	// Start with an empty list of items
	m := list.NewModel([]list.Item{}, repoListItemDelegate(cloneDirPath), 0, 0)
	m.StatusMessageLifetime = time.Second * 60
	m.Title = fmt.Sprintf("[Repositories] [%s] [%s]", strings.Join(orgs, " "), cloneDirPath)
	m.SetSpinner(spinner.MiniDot)
	m.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("H"),
				key.WithHelp("H", "toggle help"),
			),
			key.NewBinding(
				key.WithKeys("R"),
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
