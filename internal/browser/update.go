package browser

import (
	"context"
	"fmt"
	"log"
	"ogit/internal/gitutils"
	"ogit/service"
	"ogit/upstream"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("Updating UI")

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "r":
				return m, func() tea.Msg { return refreshReposMsg{} }
			default:
				log.Println("Key Pressed", string(msg.Runes))
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	return m, cmd
}

func delegateWithUpdateFunc(cloneDirPath string, orgs []string) list.DefaultDelegate {
	updateFunc := func(msg tea.Msg, m *list.Model) tea.Cmd {
		log.Println("Updating Items")

		selected, ok := m.SelectedItem().(repoListItem)
		if !ok && len(m.VisibleItems()) > 0 {
			return m.NewStatusMessage("unknown item type")
		}

		switch msg := msg.(type) {
		case refreshReposMsg:
			return tea.Batch(
				m.StartSpinner(),
				func() tea.Msg {
					s := service.NewRepositoryService(upstream.NewGithubClient(github.NewClient(nil)))
					repos, err := s.GetRepositoriesByOwners(context.Background(), orgs)
					if err != nil {
						log.Println(err)
						return updateStatusMsg(statusError(err.Error()))
					}

					return refreshReposDoneMsg{repos: *repos}
				},
			)

		case refreshReposDoneMsg:
			repos := msg.repos
			newItems := make([]list.Item, len(repos))

			for i := range repos {
				repoItem := repoListItem{
					title:       repos[i].Owner + "/" + repos[i].Name,
					owner:       repos[i].Owner,
					name:        repos[i].Name,
					description: repos[i].Description,
					browserURL:  repos[i].BrowserURL,
					cloneURL:    repos[i].CloneURL,
				}

				if repoItem.Cloned(cloneDirPath) {
					repoItem.title = statusMessageStyle(repoItem.Title())
					repoItem.description = statusMessageStyle(repoItem.Description())
				}
				newItems[i] = repoItem
			}

			m.SetItems(newItems)
			m.StopSpinner()
			return m.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Fetched %d repos", len(newItems))))

		case updateStatusMsg:
			m.StopSpinner()
			return m.NewStatusMessage(string(msg))

		case tea.KeyMsg:
			switch msg.String() {
			case "c":
				return tea.Batch(
					m.StartSpinner(),
					func() tea.Msg {
						clonePath := path.Join(cloneDirPath, selected.Owner(), selected.Name())
						if gitutils.Cloned(clonePath) {
							return updateStatusMsg(statusMessageStyle("Already Cloned"))
						}

						repoOnDisk, err := gitutils.CloneToDisk(context.Background(),
							selected.CloneURL(),
							clonePath,
							log.Default().Writer(),
						)
						if err != nil {
							return updateStatusMsg(statusError(err.Error()))
						}

						selected.title = statusMessageStyle(selected.title)
						selected.description = statusMessageStyle(selected.description)

						m.SetItem(m.Index(), selected)
						return updateStatusMsg(statusMessageStyle(repoOnDisk.String()))
					},
				)

			default:
				lastCommit, err := selected.LastCommitInfo(cloneDirPath)
				if err != nil {
					return m.NewStatusMessage(fmt.Sprintf("unable to read last commit: %s", err))
				}

				return m.NewStatusMessage(lastCommit)
			}
		}

		return nil
	}

	d := list.NewDefaultDelegate()
	d.UpdateFunc = updateFunc
	return d
}
