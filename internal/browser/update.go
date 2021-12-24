package browser

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func delegateUpdateFunc(binding key.Binding) func(msg tea.Msg, m *list.Model) tea.Cmd {

	statusMessageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
		Render

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
