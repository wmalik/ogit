package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ogit/service"
)

var _ = Describe("Repository service", func() {
	Context("When no owner is provided", func() {
		var repoService service.RepositoryService
		var repositories service.Repositories
		BeforeEach(func() {
			repositories = repoService.GetRepositoriesByOwners([]string{})
		})
		It("Returns no repository", func() {
			Expect(repositories).To(Equal(service.Repositories{}))
		})
	})
})
