package main

import (
	"log"
	"os"

	"ogit/internal/browser"
	"ogit/internal/gitconfig"
	"ogit/service"
	"ogit/upstream"

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

	rs := service.NewRepositoryService(upstream.NewGithubClientWithToken(os.Getenv("GITHUB_TOKEN")))
	if err := tea.NewProgram(
		browser.NewModel(gitConf.Orgs(), gitConf.CloneDirPath(), rs),
	).Start(); err != nil {
		log.Fatalln(err)
	}
}
