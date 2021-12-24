package main

import (
	"log"

	"ogit/internal/browser"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	f, err := tea.LogToFile("/tmp/ogit.log", "debug")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	if err := tea.NewProgram(browser.NewModel()).Start(); err != nil {
		log.Fatalln(err)
	}
}
