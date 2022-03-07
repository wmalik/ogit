package browser

import (
	"context"
	"log"
	"ogit/internal/gitutils"
	"ogit/internal/utils"

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
		m.list.StopSpinner()
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
func delegateItemUpdate(storagePath string, gu *gitutils.GitUtils) list.DefaultDelegate {
	updateFunc := func(msg tea.Msg, m *list.Model) tea.Cmd {
		log.Println("Updating Item")

		selected, ok := m.SelectedItem().(repoItem)
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
				if selected.Cloned() {
					return updateStatusMsg(statusMessageStyle("Already Cloned"))
				}

				repoString, err := selected.Clone(context.Background(), gu)
				if err != nil {
					return updateStatusMsg(statusError(err.Error()))
				}

				selected.SetTitle(brightStyle.Render(selected.Repository.Title))

				m.SetItem(m.Index(), selected)
				return updateStatusMsg(statusMessageStyle(repoString))
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
					return openURLMsg(selected.Repository.BrowserHomepageURL)
				}

			case "p":
				return func() tea.Msg {
					return openURLMsg(selected.Repository.BrowserPullRequestsURL)
				}

			default:
				return m.NewStatusMessage(selected.Repository.Description)
			}
		}

		return nil
	}

	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(dimmedColor)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.UnsetForeground().Background(selectedColor)
	d.ShowDescription = false
	d.SetSpacing(0)
	d.UpdateFunc = updateFunc
	return d
}
