package main

import (
	"log"

	"ogit/internal/browser"
	"ogit/internal/gitconfig"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := tea.LogToFile("/tmp/ogit.log", "debug")
	if err != nil {
		log.Fatalln(err)
	}

	defer f.Close()
	if err := tea.NewProgram(browser.NewModel(gitConf.Orgs(), gitConf.CloneDirPath())).Start(); err != nil {
		log.Fatalln(err)
	}
}
