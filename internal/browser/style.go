package browser

import "github.com/charmbracelet/lipgloss"

var appStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingLeft(5).PaddingTop(5).PaddingBottom(5).PaddingRight(5)

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#0000ff", Dark: "#000099"})

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#25A065")).
	Padding(0, 1)

func statusMessageStyle(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
		Render(str)
}

func statusError(str string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#eb4f34", Dark: "#eb4f34"}).
		Render(str)
}
