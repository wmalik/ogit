package upstream_test

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ogit/mock"
	"ogit/upstream"
)

var _ = Describe("Github repo", func() {
	var client *upstream.GithubClient
	var repositories []upstream.HostRepository
	BeforeEach(func() {
		httpClient := mock.NewHTTPClient().
			Mock("GET", "/users/greatuser/repos",
				func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(`[
						{
							"name": "dotfiles",
							"full_name": "greatuser/dotfiles",
							"private": false,
							"owner": {
								"login": "greatuser"
							}
						},
						{
							"name": "personal-website",
							"full_name": "greatuser/personal-website",
							"private": false,
							"owner": {
								"login": "greatuser"
							}
						}
					]`))
				},
			).Client()
		client = upstream.NewGithubClient(github.NewClient(httpClient))
		repositories = client.GetRepositories(context.Background(), []string{"greatuser"})
	})
	It("Returns the matching repositories", func() {
		Expect(len(repositories)).To(Equal(2))
		Expect(repositories[0].GetName()).To(Equal("dotfiles"))
		Expect(repositories[1].GetName()).To(Equal("personal-website"))
	})
})
