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
		var repositories *service.Repositories
		var err error
		BeforeEach(func() {
			repoService = service.NewRepositoryService(upstream.NewMockClient(), false)
			repositories, err = repoService.GetRepositoriesByOwners(context.Background(), []string{})
			Expect(err).To(BeNil())
		})
		It("Returns no repository", func() {
			Expect(*repositories).To(Equal(service.Repositories{}))
		})
	})
	Context("When an owner is provided", func() {
		var repoService *service.RepositoryService
		var repositories *service.Repositories
		var err error
		BeforeEach(func() {
			client := upstream.NewMockClient().WithRepositories([]upstream.MockRepository{
				{
					Owner:                  "wmalik",
					Name:                   "ogit",
					Description:            "TUI for browsing GitHub and GitLab orgnizations",
					BrowserHomepageURL:     "https://github.com/wmalik/ogit",
					BrowserPullRequestsURL: "https://github.com/wmalik/ogit/pulls",
					CloneURL:               "https://github.com/wmalik/ogit.git",
				},
				{
					Owner:                  "wmalik",
					Name:                   "dotfiles",
					Description:            "wmalik's config files",
					BrowserHomepageURL:     "https://github.com/wmalik/dotfiles",
					BrowserPullRequestsURL: "https://github.com/wmalik/dotfiles/pulls",
					CloneURL:               "https://github.com/wmalik/dotfiles.git",
				},
				{
					Owner:                  "padawin",
					Name:                   "dotfiles",
					Description:            "padawin's config files",
					BrowserHomepageURL:     "https://github.com/padawin/dotfiles",
					BrowserPullRequestsURL: "https://github.com/padawin/dotfiles/pulls",
					CloneURL:               "https://github.com/padawin/dotfiles.git",
				},
			})
			repoService = service.NewRepositoryService(client, false)
			repositories, err = repoService.GetRepositoriesByOwners(context.Background(), []string{"wmalik"})
			Expect(err).To(BeNil())
		})
		It("Returns the matching repositories", func() {
			Expect(len(*repositories)).To(Equal(2))
			Expect((*repositories)[0].Name).To(Equal("ogit"))
			Expect((*repositories)[0].Description).To(Equal("TUI for browsing GitHub and GitLab orgnizations"))
			Expect((*repositories)[0].BrowserHomepageURL).To(Equal("https://github.com/wmalik/ogit"))
			Expect((*repositories)[0].BrowserPullRequestsURL).To(Equal("https://github.com/wmalik/ogit/pulls"))
			Expect((*repositories)[0].CloneURL).To(Equal("https://github.com/wmalik/ogit.git"))
			Expect((*repositories)[1].Name).To(Equal("dotfiles"))
			Expect((*repositories)[1].Description).To(Equal("wmalik's config files"))
			Expect((*repositories)[1].BrowserHomepageURL).To(Equal("https://github.com/wmalik/dotfiles"))
			Expect((*repositories)[1].BrowserPullRequestsURL).To(Equal("https://github.com/wmalik/dotfiles/pulls"))
			Expect((*repositories)[1].CloneURL).To(Equal("https://github.com/wmalik/dotfiles.git"))
		})
	})
})
