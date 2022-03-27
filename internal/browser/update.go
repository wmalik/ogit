package browser

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/wmalik/ogit/internal/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	selected, ok := m.list.SelectedItem().(repoItem)
	if !ok && len(m.list.VisibleItems()) > 0 {
		return m, nil
	}

	if m.list.FilterState() != list.Filtering {
		cmds = append(cmds, handleKeyMsg(msg, m, selected))
		cmds = append(cmds, handleMsg(msg, m))
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	m.selectedItemStoragePath = selected.repoStoragePath
	return m, tea.Batch(append(cmds, cmd)...)
}

func handleMsg(msg tea.Msg, m *Model) tea.Cmd {
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
			if msg.repo.Cloned() {
				return updateBottomStatusBarMsg(statusMessageStyle("[Already Cloned] " + msg.repo.StoragePath()))
			}

			repoString, err := msg.repo.Clone(context.Background(), m.gu)
			if err != nil {
				return updateBottomStatusBarMsg(statusError(err.Error()))
			}

			msg.repo.SetTitle(clonedRepoStyle.Render(msg.repo.Repository.Title))

			m.list.SetItem(msg.index, msg.repo)
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
	return tea.Batch(cmds...)
}

func handleKeyMsg(msg tea.Msg, m *Model, selected repoItem) tea.Cmd {
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "o":
			if !selected.Cloned() {
				return func() tea.Msg {
					return updateBottomStatusBarMsg(
						statusError("Not cloned yet, press c to clone"),
					)
				}
			}
			rangerCmd := exec.Command("xdg-open", selected.StoragePath()) //nolint:gosec
			rangerCmd.Stdin = os.Stdin
			rangerCmd.Stdout = os.Stdout
			if err := rangerCmd.Run(); err != nil {
				return func() tea.Msg {
					return updateBottomStatusBarMsg(
						statusError(fmt.Sprintf("Unable to run xdg-open: %s", err)),
					)
				}
			}
			return tea.HideCursor

		case "v":
			if !selected.Cloned() {
				return func() tea.Msg {
					return updateBottomStatusBarMsg(
						statusError("Not cloned yet, press c to clone"),
					)
				}
			}
			vimCmd := exec.Command("vim", selected.StoragePath()) //nolint:gosec
			vimCmd.Stdin = os.Stdin
			vimCmd.Stdout = os.Stdout
			if err := vimCmd.Run(); err != nil {
				return func() tea.Msg {
					return updateBottomStatusBarMsg(
						statusError(fmt.Sprintf("Unable to run vim: %s", err)),
					)
				}
			}
			return tea.HideCursor

		case "enter":
			if !selected.Cloned() {
				return func() tea.Msg {
					return updateBottomStatusBarMsg(
						statusError("Not cloned yet, press c to clone"),
					)
				}
			}
			m.spawnShell = true
			m.shellArgs = []string{"-i"}
			cmds = append(cmds, tea.Quit)
		case "c":
			cmds = append(cmds, tea.Batch(
				m.list.StartSpinner(),
				func() tea.Msg {
					return cloneRepoMsg{selected, m.list.Index()}
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

	return tea.Batch(cmds...)
}

// listItemDelegate configures general behaviour/styling of the list items.
func listItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(dimmedColor)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Background(selectedColorBg).Foreground(selectedColorFg)
	d.ShowDescription = false
	d.SetSpacing(0)
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		if selected, ok := m.SelectedItem().(repoItem); ok {
			return m.NewStatusMessage(selected.Description())
		}
		return nil
	}
	return d
}
