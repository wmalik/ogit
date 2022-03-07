package service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/wmalik/ogit/service"
	"github.com/wmalik/ogit/upstream"
)

var _ = Describe("Repository service", func() {
	Context("When no owner is provided", func() {
		var repoService *service.RepositoryService
		var repositories *service.Repositories
		var err error
		BeforeEach(func() {
			gitlabClient := upstream.NewMockClient()
			repoService = service.NewRepositoryService(upstream.NewMockClient(), gitlabClient, false)
			repositories, err = repoService.GetRepositoriesByOwners(context.Background(), []string{}, []string{})
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
			gitlabClient := upstream.NewMockClient().WithRepositories([]upstream.MockRepository{
				{
					Provider:               "gitlab",
					Owner:                  "wmalik",
					Name:                   "ogit",
					Description:            "TUI for browsing GitHub and GitLab orgnizations",
					BrowserHomepageURL:     "https://gitlab.com/wmalik/ogit",
					BrowserPullRequestsURL: "https://gitlab.com/wmalik/ogit/pulls",
					HTTPSCloneURL:          "https://gitlab.com/wmalik/ogit.git",
					SSHCloneURL:            "git@gitlab.com/wmalik/ogit.git",
				},
				{
					Provider:               "gitlab",
					Owner:                  "wmalik",
					Name:                   "dotfiles",
					Description:            "wmalik's config files",
					BrowserHomepageURL:     "https://gitlab.com/wmalik/dotfiles",
					BrowserPullRequestsURL: "https://gitlab.com/wmalik/dotfiles/pulls",
					HTTPSCloneURL:          "https://gitlab.com/wmalik/dotfiles.git",
					SSHCloneURL:            "git@gitlab.com/wmalik/dotfiles.git",
				},
				{
					Provider:               "gitlab",
					Owner:                  "padawin",
					Name:                   "dotfiles",
					Description:            "padawin's config files",
					BrowserHomepageURL:     "https://gitlab.com/padawin/dotfiles",
					BrowserPullRequestsURL: "https://gitlab.com/padawin/dotfiles/pulls",
					HTTPSCloneURL:          "https://gitlab.com/padawin/dotfiles.git",
					SSHCloneURL:            "git@gitlab.com/padawin/dotfiles.git",
				},
			})
			client := upstream.NewMockClient().WithRepositories([]upstream.MockRepository{
				{
					Provider:               "github",
					Owner:                  "wmalik",
					Name:                   "ogit",
					Description:            "TUI for browsing GitHub and GitLab orgnizations",
					BrowserHomepageURL:     "https://github.com/wmalik/ogit",
					BrowserPullRequestsURL: "https://github.com/wmalik/ogit/pulls",
					HTTPSCloneURL:          "https://github.com/wmalik/ogit.git",
					SSHCloneURL:            "git@github.com/wmalik/ogit.git",
				},
				{
					Provider:               "github",
					Owner:                  "wmalik",
					Name:                   "dotfiles",
					Description:            "wmalik's config files",
					BrowserHomepageURL:     "https://github.com/wmalik/dotfiles",
					BrowserPullRequestsURL: "https://github.com/wmalik/dotfiles/pulls",
					HTTPSCloneURL:          "https://github.com/wmalik/dotfiles.git",
					SSHCloneURL:            "git@github.com/wmalik/dotfiles.git",
				},
				{
					Provider:               "github",
					Owner:                  "padawin",
					Name:                   "dotfiles",
					Description:            "padawin's config files",
					BrowserHomepageURL:     "https://github.com/padawin/dotfiles",
					BrowserPullRequestsURL: "https://github.com/padawin/dotfiles/pulls",
					HTTPSCloneURL:          "https://github.com/padawin/dotfiles.git",
					SSHCloneURL:            "git@github.com/padawin/dotfiles.git",
				},
			})
			repoService = service.NewRepositoryService(client, gitlabClient, false)
			repositories, err = repoService.GetRepositoriesByOwners(context.Background(), []string{"wmalik"}, []string{"wmalik"})
			Expect(err).To(BeNil())
		})
		It("Returns the matching repositories", func() {
			Expect(len(*repositories)).To(Equal(4))
			Expect((*repositories)[0].Provider).To(Equal("github"))
			Expect((*repositories)[0].Name).To(Equal("ogit"))
			Expect((*repositories)[0].Description).To(Equal("TUI for browsing GitHub and GitLab orgnizations"))
			Expect((*repositories)[0].BrowserHomepageURL).To(Equal("https://github.com/wmalik/ogit"))
			Expect((*repositories)[0].BrowserPullRequestsURL).To(Equal("https://github.com/wmalik/ogit/pulls"))
			Expect((*repositories)[0].HTTPSCloneURL).To(Equal("https://github.com/wmalik/ogit.git"))
			Expect((*repositories)[0].SSHCloneURL).To(Equal("git@github.com/wmalik/ogit.git"))
			Expect((*repositories)[1].Provider).To(Equal("github"))
			Expect((*repositories)[1].Name).To(Equal("dotfiles"))
			Expect((*repositories)[1].Description).To(Equal("wmalik's config files"))
			Expect((*repositories)[1].BrowserHomepageURL).To(Equal("https://github.com/wmalik/dotfiles"))
			Expect((*repositories)[1].BrowserPullRequestsURL).To(Equal("https://github.com/wmalik/dotfiles/pulls"))
			Expect((*repositories)[1].HTTPSCloneURL).To(Equal("https://github.com/wmalik/dotfiles.git"))
			Expect((*repositories)[1].SSHCloneURL).To(Equal("git@github.com/wmalik/dotfiles.git"))
			Expect((*repositories)[2].Provider).To(Equal("gitlab"))
			Expect((*repositories)[2].Name).To(Equal("ogit"))
			Expect((*repositories)[2].Description).To(Equal("TUI for browsing GitHub and GitLab orgnizations"))
			Expect((*repositories)[2].BrowserHomepageURL).To(Equal("https://gitlab.com/wmalik/ogit"))
			Expect((*repositories)[2].BrowserPullRequestsURL).To(Equal("https://gitlab.com/wmalik/ogit/pulls"))
			Expect((*repositories)[2].HTTPSCloneURL).To(Equal("https://gitlab.com/wmalik/ogit.git"))
			Expect((*repositories)[2].SSHCloneURL).To(Equal("git@gitlab.com/wmalik/ogit.git"))
			Expect((*repositories)[3].Provider).To(Equal("gitlab"))
			Expect((*repositories)[3].Name).To(Equal("dotfiles"))
			Expect((*repositories)[3].Description).To(Equal("wmalik's config files"))
			Expect((*repositories)[3].BrowserHomepageURL).To(Equal("https://gitlab.com/wmalik/dotfiles"))
			Expect((*repositories)[3].BrowserPullRequestsURL).To(Equal("https://gitlab.com/wmalik/dotfiles/pulls"))
			Expect((*repositories)[3].HTTPSCloneURL).To(Equal("https://gitlab.com/wmalik/dotfiles.git"))
			Expect((*repositories)[3].SSHCloneURL).To(Equal("git@gitlab.com/wmalik/dotfiles.git"))
		})
	})
})
