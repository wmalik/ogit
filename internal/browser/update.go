package browser

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ogit/internal/gitutils"
	"ogit/service"
	"ogit/upstream"
	"path"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
)

type doneFetchRepoList struct {
	repos service.Repositories
}

type errorFetchRepoList struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errorFetchRepoList) Error() string { return e.err.Error() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("Updating UI")
	var cmds []tea.Cmd

	// TODO: move this to a key mapping handler
	if m.fetch {
		m.fetch = false
		cmds = append(cmds, m.list.StartSpinner())

		fetchReposListCmd := func() tea.Msg {
			s := service.NewRepositoryService(upstream.NewGithubClient(github.NewClient(nil)))
			repos, err := s.GetRepositoriesByOwners(context.Background(), m.orgs)
			if err != nil {
				log.Println(err)
				return errorFetchRepoList{err}
			}

			return doneFetchRepoList{
				repos: *repos,
			}
		}

		cmds = append(cmds, fetchReposListCmd)
	}

	switch msg := msg.(type) {
	case doneFetchRepoList:
		repos := doneFetchRepoList(msg).repos
		repoListItems := make([]list.Item, len(repos))
		for i := 0; i < len(repos); i++ {
			repoItem := repoListItem{
				title:       repos[i].Owner + "/" + repos[i].Name,
				owner:       repos[i].Owner,
				name:        repos[i].Name,
				description: repos[i].Description,
				browserURL:  repos[i].BrowserURL,
				cloneURL:    repos[i].CloneURL,
			}

			if repoItem.Cloned(m.cloneDirPath) {
				repoItem.title = statusMessageStyle(repoItem.Title())
				repoItem.description = statusMessageStyle(repoItem.Description())
			}
			repoListItems[i] = repoItem
		}

		m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Fetched %d repos", len(repoListItems))))
		m.list.SetItems(repoListItems)
		m.list.StopSpinner()

	case errorFetchRepoList:
		m.list.StopSpinner()
		m.list.StatusMessageLifetime = time.Second * 10
		cmds = append(cmds, m.list.NewStatusMessage(statusError(errorFetchRepoList(msg).Error())))

	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "H":
				m.list.SetShowHelp(!m.list.ShowHelp())
				return m, nil
			case "R":
				m.fetch = true
				return m, func() tea.Msg { return nil }
			default:
				log.Println("Key Pressed", string(msg.Runes))
			}
		}
	}

	// This will also call the delegate's update function
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func delegateUpdateFunc(binding key.Binding, cloneDirPath string) func(msg tea.Msg, m *list.Model) tea.Cmd {
	return func(msg tea.Msg, m *list.Model) tea.Cmd {
		log.Println("Updating delegate UI")

		var title, owner, name, browserURL, cloneURL string

		selectedItem := m.SelectedItem()
		selectedRepoListItem, ok := selectedItem.(repoListItem)
		// TODO: if !ok return nil
		if ok {
			title = selectedRepoListItem.Title()
			owner = selectedRepoListItem.Owner()
			name = selectedRepoListItem.Name()
			browserURL = selectedRepoListItem.BrowserURL()
			cloneURL = selectedRepoListItem.CloneURL()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case doneCloneToDisk:
			m.StopSpinner()

			done := doneCloneToDisk(msg)
			clonedItem := done.item
			clonedItemIndex := done.index
			clonedItem.title = statusMessageStyle(clonedItem.Title())
			clonedItem.description = statusMessageStyle(clonedItem.Description())
			log.Println(done.repo.LastCommit())

			return tea.Batch(
				m.NewStatusMessage(statusMessageStyle("Cloned "+done.repo.String())),
				m.SetItem(clonedItemIndex, clonedItem),
			)

		case errCloneToDisk:
			m.StopSpinner()
			if errors.Is(errCloneToDisk(msg).err, gitutils.ErrRepoAlreadyCloned) {
				return m.NewStatusMessage(statusError(errCloneToDisk(msg).err.Error()))
			}
			return m.NewStatusMessage(statusError(errCloneToDisk(msg).err.Error()))

		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return m.NewStatusMessage(statusMessageStyle("You chose " + title))
			case tea.KeyRunes:
				switch string(msg.Runes) {
				case "o":
					return m.NewStatusMessage(statusMessageStyle("Opening in firefox: " + browserURL))
				case "c":

					cloneToDiskCmd := func() tea.Msg {
						progress := log.Default().Writer()
						diskPath := path.Join(cloneDirPath, owner, name)
						repoOnDisk, err := gitutils.CloneToDisk(context.Background(), cloneURL, diskPath, progress)
						if err != nil {
							log.Println(err)
							return errCloneToDisk{err: err}
						}

						return doneCloneToDisk{index: m.Index(), item: selectedRepoListItem, repo: repoOnDisk}
					}
					return tea.Batch(
						m.StartSpinner(),
						cloneToDiskCmd,
						m.SetItem(m.Index(), selectedRepoListItem),
					)

				}
			}
		}

		return nil
	}
}

type errCloneToDisk struct {
	err error
}

type doneCloneToDisk struct {
	index int
	item  repoListItem
	repo  *gitutils.Repository
}
