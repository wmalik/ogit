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
				{Owner: "wmalik", Name: "ogit"},
				{Owner: "wmalik", Name: "dotfiles"},
				{Owner: "padawin", Name: "dotfiles"},
			})
			repoService = service.NewRepositoryService(client)
			repositories = repoService.GetRepositoriesByOwners(context.Background(), []string{"wmalik"})
		})
		It("Returns the matching repositories", func() {
			Expect(len(repositories)).To(Equal(2))
			Expect(repositories[0].Name).To(Equal("ogit"))
			Expect(repositories[1].Name).To(Equal("dotfiles"))
		})
	})
})
