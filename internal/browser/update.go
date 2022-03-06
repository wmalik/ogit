package browser

import (
	"context"
	"log"
	"ogit/internal/gitutils"
	"ogit/internal/utils"
	"path"

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
			key.WithKeys("w"),
			key.WithHelp("w", "browse home page"),
		),
		key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "browse pull requests"),
		),
	}
}

// Update is called whenever the whole model is updated
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("Updating UI")

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		bottomGap = bottomGap + bottomStatusBarStyle.GetHeight()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case updateBottomStatusBarMsg:
		m.bottomStatusBar = string(msg)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			default:
				log.Println("Key Pressed", string(msg.Runes))
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

// delegateItemUpdate is called whenever a specific item is updated.
// It is used for example for messages like "clone repo"
func delegateItemUpdate(cloneDirPath string, gu *gitutils.GitUtils) list.DefaultDelegate {
	updateFunc := func(msg tea.Msg, m *list.Model) tea.Cmd {
		log.Println("Updating Item")

		selected, ok := m.SelectedItem().(repoListItem)
		if !ok && len(m.VisibleItems()) > 0 {
			return m.NewStatusMessage("unknown item type")
		}

		switch msg := msg.(type) {
		case updateStatusMsg:
			m.StopSpinner()
			return m.NewStatusMessage(string(msg))

		case cloneRepoMsg:
			return func() tea.Msg {
				defer m.StopSpinner()
				clonePath := path.Join(cloneDirPath, selected.Owner(), selected.Name())
				if gitutils.Cloned(clonePath) {
					return updateStatusMsg(statusMessageStyle("Already Cloned"))
				}

				repoOnDisk, err := gu.CloneToDisk(context.Background(),
					selected.HTTPSCloneURL(),
					selected.SSHCloneURL(),
					clonePath,
					log.Default().Writer(),
				)
				if err != nil {
					return updateStatusMsg(statusError(err.Error()))
				}

				selected.title = brightStyle.Render(selected.title)

				m.SetItem(m.Index(), selected)
				return updateStatusMsg(statusMessageStyle(repoOnDisk.String()))
			}

		case openURLMsg:
			return func() tea.Msg {
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
			}

		case tea.KeyMsg:
			switch msg.String() {
			case "c":
				return tea.Batch(
					m.StartSpinner(),
					func() tea.Msg {
						return cloneRepoMsg{selected}
					},
				)
			case "w":
				return func() tea.Msg {
					return openURLMsg(selected.BrowserHomepageURL())
				}

			case "p":
				return func() tea.Msg {
					return openURLMsg(selected.BrowserPullRequestsURL())
				}

			default:
				return m.NewStatusMessage(selected.description)
			}
		}

		return nil
	}

	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.UpdateFunc = updateFunc
	return d
}
