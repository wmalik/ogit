package browser

import tea "github.com/charmbracelet/bubbletea"

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return refreshReposMsg{}
	}
}
