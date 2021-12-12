package main

import (
	"context"
	"fmt"
	"log"
	"ogit/service"
	"ogit/upstream"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
	"github.com/tcnksm/go-gitconfig"
)

func main() {
	orgs, err := orgsFromGitConfig()
	if err != nil {
		log.Fatalln(err)
	}

	s := service.NewRepositoryService(upstream.NewGithubClient(github.NewClient(nil)))
	repos := s.GetRepositoriesByOwners(context.Background(), orgs)

	if err := tea.NewProgram(newModel(repos)).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// orgsFromGitConfig loads the value of ogit.orgs from ~/.gitconfig
func orgsFromGitConfig() ([]string, error) {
	orgsRaw, err := gitconfig.Entire("ogit.orgs")
	if err != nil {
		return nil, fmt.Errorf("unable to read git config: %s", err)
	}

	if orgsRaw == "" {
		return nil, fmt.Errorf("Please configure ogit.orgs in your ~/.gitconfig")
	}

	orgs := []string{}
	for _, org := range strings.Split(orgsRaw, ",") {
		orgs = append(orgs, strings.TrimSpace(org))
	}

	return orgs, nil
}
