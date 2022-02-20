package main

import (
	"log"
	"os"

	"ogit/internal/browser"
	"ogit/internal/gitconfig"
	"ogit/internal/gitutils"
	"ogit/service"
	"ogit/upstream"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	gitConf, err := gitconfig.ReadGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	gu, err := gitutils.NewGitUtils(gitConf.UseSSHAgent(), gitConf.PrivKeyPath())
	if err != nil {
		log.Fatalln(err)
	}

	f, err := tea.LogToFile("/tmp/ogit.log", "debug")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	gitlabClient, err := upstream.NewGitlabClientWithToken(os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	rs := service.NewRepositoryService(
		upstream.NewGithubClientWithToken(os.Getenv("GITHUB_TOKEN")),
		gitlabClient,
		gitConf.FetchAuthenticatedUserRepos(),
	)
	if err := tea.NewProgram(
		browser.NewModel(gitConf.Orgs(), gitConf.GitlabGroups(), gitConf.CloneDirPath(), rs, gu),
	).Start(); err != nil {
		log.Fatalln(err)
	}
}
