package browser

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		appStyle.Render(m.list.View()),
		bottomStatusBarStyle.Render(m.bottomStatusBar),
	)
}
