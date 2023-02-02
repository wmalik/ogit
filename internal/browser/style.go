package browser

import (
	"github.com/charmbracelet/lipgloss"
)

var appStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingLeft(2).
	PaddingTop(2).
	PaddingBottom(2).
	PaddingRight(2)

var bottomStatusBarStyle = lipgloss.NewStyle().Height(1).Faint(true)

var clonedRepoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#030303", Dark: "#dddddd"})

var defaultDimmedColorFg = lipgloss.AdaptiveColor{Light: "#79787a", Dark: "#7F7C82"}
var defaultSelectedColorBg = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#333333"}
var defaultSelectedColorFg = lipgloss.AdaptiveColor{Light: "#79787a", Dark: "#7F7C82"}
var defaultTitleFg = lipgloss.AdaptiveColor{Light: "#f5f2f6", Dark: "#f5f2f6"}
var defaultTitleBg = lipgloss.AdaptiveColor{Light: "#5b186e", Dark: "#5b186e"}

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

func getColor(color *lipgloss.AdaptiveColor, defaultColor lipgloss.AdaptiveColor) lipgloss.AdaptiveColor {
	if color == nil {
		return defaultColor
	}
	return *color
}
