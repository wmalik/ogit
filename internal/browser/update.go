package browser

import (
	"context"
	"log"

	"github.com/wmalik/ogit/internal/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func availableKeyBindingsCB() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "clone a repository (shallow)"),
		),
		key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open cloned directory"),
		),
		key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "browse home page"),
		),
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "browse pull requests"),
		),
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	selected, ok := m.list.SelectedItem().(repoItem)
	if !ok && len(m.list.VisibleItems()) > 0 {
		return m, nil
	}

	if m.list.FilterState() != list.Filtering {
		cmds = append(cmds, handleKeyMsg(msg, m, selected)...)
		cmds = append(cmds, handleMsg(msg, m, selected)...)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	m.selectedItemStoragePath = selected.repoStoragePath
	return m, tea.Batch(append(cmds, cmd)...)
}

func handleMsg(msg tea.Msg, m *model, selected repoItem) []tea.Cmd {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		bottomGap = bottomGap + bottomStatusBarStyle.GetHeight()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case updateBottomStatusBarMsg:
		m.list.StopSpinner()
		m.bottomStatusBar = string(msg)

	case updateStatusMsg:
		m.list.StopSpinner()
		cmds = append(cmds, m.list.NewStatusMessage(string(msg)))

	case cloneRepoMsg:
		cmds = append(cmds, func() tea.Msg {
			defer m.list.StopSpinner()
			if selected.Cloned() {
				return updateBottomStatusBarMsg(statusMessageStyle("[Already Cloned] " + selected.StoragePath()))
			}

			repoString, err := selected.Clone(context.Background(), m.gu)
			if err != nil {
				return updateBottomStatusBarMsg(statusError(err.Error()))
			}

			selected.SetTitle(brightStyle.Render(selected.Repository.Title))

			m.list.SetItem(m.list.Index(), selected)
			return updateBottomStatusBarMsg(statusMessageStyle("[Cloned] " + repoString))
		})

	case openURLMsg:
		cmds = append(cmds, func() tea.Msg {
			u := string(msg)
			if u == "" {
				return updateStatusMsg(statusError("URL not available"))
			}
			err := utils.OpenURL(u)
			if err != nil {
				log.Println(err)
				return updateStatusMsg(statusError(err.Error()))
			}
			return nil
		})

	}
	return cmds
}

func handleKeyMsg(msg tea.Msg, m *model, selected repoItem) []tea.Cmd {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch string(msg.Runes) {
		case "o":
			if selected.Cloned() {
				m.spawnShell = true
				cmds = append(cmds, tea.Quit)
				break
			}
			cmds = append(cmds, func() tea.Msg {
				return updateBottomStatusBarMsg(
					statusError("Not cloned yet, press c to clone"),
				)
			})
		case "c":
			cmds = append(cmds, tea.Batch(
				m.list.StartSpinner(),
				func() tea.Msg {
					return cloneRepoMsg{selected}
				},
			))
		case "w":
			cmds = append(cmds, func() tea.Msg {
				return openURLMsg(selected.Repository.BrowserHomepageURL)
			})
		case "p":
			cmds = append(cmds, func() tea.Msg {
				return openURLMsg(selected.Repository.BrowserPullRequestsURL)
			})
		default:
			log.Println("Key Pressed", string(msg.Runes))
		}
	}

	return cmds
}

// listItemDelegate configures general behaviour/styling of the list items
func listItemDelegate(storagePath string) list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(dimmedColor)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.UnsetForeground().Background(selectedColor)
	d.ShowDescription = false
	d.SetSpacing(0)
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		if selected, ok := m.SelectedItem().(repoItem); ok {
			return m.NewStatusMessage(selected.Repository.Description)
		}
		return nil
	}
	return d
}
