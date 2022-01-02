package browser

import tea "github.com/charmbracelet/bubbletea"

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return fetchAPIUsageMsg{} },
		func() tea.Msg { return refreshReposMsg{} },
	)
}
