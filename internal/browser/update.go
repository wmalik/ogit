package browser

import (
	"context"
	"fmt"
	"log"
	"ogit/service"
	"ogit/upstream"

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
	var cmds []tea.Cmd

	if m.fetch {
		m.fetch = false
		cmds = append(cmds, m.list.StartSpinner())
		log.Println("inside fetch")
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
			log.Println("Adding repo", repos[i].Name)
			repoListItems[i] = repoListItem{
				title:       repos[i].Owner + "/" + repos[i].Name,
				description: repos[i].Description,
				browserURL:  repos[i].BrowserURL,
				cloneURL:    repos[i].CloneURL,
			}
		}

		m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Fetched %d repos", len(repoListItems))))
		m.list.SetItems(repoListItems)
		m.list.StopSpinner()

	case errorFetchRepoList:
		m.list.StopSpinner()
		cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle(errorFetchRepoList(msg).Error())))

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

	log.Println("finishing update")
	// This will also call the delegate's update function
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func delegateUpdateFunc(binding key.Binding) func(msg tea.Msg, m *list.Model) tea.Cmd {

	return func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title, browserURL, cloneURL string

		if i, ok := m.SelectedItem().(repoListItem); ok {
			title = i.Title()
			browserURL = i.BrowserURL()
			cloneURL = i.CloneURL()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				return m.NewStatusMessage(statusMessageStyle("You chose " + title))
			case tea.KeyRunes:
				switch string(msg.Runes) {
				case "o":
					return m.NewStatusMessage(statusMessageStyle("Opening in firefox: " + browserURL))
				case "c":
					return m.NewStatusMessage(statusMessageStyle("Cloning " + cloneURL))
				}
			}
		}

		return nil
	}
}
