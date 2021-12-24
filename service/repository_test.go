package service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ogit/service"
	"ogit/upstream"
)

var _ = Describe("Repository service", func() {
	Context("When no owner is provided", func() {
		var repoService *service.RepositoryService
		var repositories service.Repositories
		BeforeEach(func() {
			repoService = service.NewRepositoryService(upstream.NewMockClient())
			repositories = repoService.GetRepositoriesByOwners(context.Background(), []string{})
		})
		It("Returns no repository", func() {
			Expect(repositories).To(Equal(service.Repositories{}))
		})
	})
	Context("When an owner is provided", func() {
		var repoService *service.RepositoryService
		var repositories service.Repositories
		BeforeEach(func() {
			client := upstream.NewMockClient().WithRepositories([]upstream.MockRepository{
				{
					Owner:       "wmalik",
					Name:        "ogit",
					Description: "TUI for browsing GitHub and GitLab orgnizations",
					BrowserURL:  "https://github.com/wmalik/ogit",
					CloneURL:    "https://github.com/wmalik/ogit.git",
				},
				{
					Owner:       "wmalik",
					Name:        "dotfiles",
					Description: "wmalik's config files",
					BrowserURL:  "https://github.com/wmalik/dotfiles",
					CloneURL:    "https://github.com/wmalik/dotfiles.git",
				},
				{
					Owner:       "padawin",
					Name:        "dotfiles",
					Description: "padawin's config files",
					BrowserURL:  "https://github.com/padawin/dotfiles",
					CloneURL:    "https://github.com/padawin/dotfiles.git",
				},
			})
			repoService = service.NewRepositoryService(client)
			repositories = repoService.GetRepositoriesByOwners(context.Background(), []string{"wmalik"})
		})
		It("Returns the matching repositories", func() {
			Expect(len(repositories)).To(Equal(2))
			Expect(repositories[0].Name).To(Equal("ogit"))
			Expect(repositories[0].Description).To(Equal("TUI for browsing GitHub and GitLab orgnizations"))
			Expect(repositories[0].BrowserURL).To(Equal("https://github.com/wmalik/ogit"))
			Expect(repositories[0].CloneURL).To(Equal("https://github.com/wmalik/ogit.git"))
			Expect(repositories[1].Name).To(Equal("dotfiles"))
			Expect(repositories[1].Description).To(Equal("wmalik's config files"))
			Expect(repositories[1].BrowserURL).To(Equal("https://github.com/wmalik/dotfiles"))
			Expect(repositories[1].CloneURL).To(Equal("https://github.com/wmalik/dotfiles.git"))
		})
	})
})
