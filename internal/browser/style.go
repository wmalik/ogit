package browser

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingLeft(5).PaddingTop(5).PaddingBottom(5).PaddingRight(5)

var bottomStatusBarStyle = lipgloss.NewStyle().Height(1).Faint(true)

var clonedRepoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#030303", Dark: "#dddddd"})

var dimmedColor = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#7F7C82"}
var selectedColorBg = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#333333"}
var selectedColorFg = lipgloss.AdaptiveColor{Light: "#79787a", Dark: "#7F7C82"}
var titleBarStyle = list.DefaultStyles().TitleBar.Background(lipgloss.Color("#52006A")).Padding(0, 1)

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
