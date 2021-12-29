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
			repoService = service.NewRepositoryService(upstream.NewMockClient())
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
				{Owner: "wmalik", Name: "ogit"},
				{Owner: "wmalik", Name: "dotfiles"},
				{Owner: "padawin", Name: "dotfiles"},
			})
			repoService = service.NewRepositoryService(client)
			repositories, err = repoService.GetRepositoriesByOwners(context.Background(), []string{"wmalik"})
			Expect(err).To(BeNil())
		})
		It("Returns the matching repositories", func() {
			Expect(len(*repositories)).To(Equal(2))
			Expect((*repositories)[0].Name).To(Equal("ogit"))
			Expect((*repositories)[1].Name).To(Equal("dotfiles"))
		})
	})
})
