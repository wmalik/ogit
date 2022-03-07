package upstream_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xanzy/go-gitlab"

	"ogit/mock"
	"ogit/upstream"
)

var _ = Describe("Gitlab repo", func() {
	var client *upstream.GitlabClient
	var gitlabClient *gitlab.Client
	var repositories []upstream.HostRepository
	var err error
	BeforeEach(func() {
		httpClient := mock.NewHTTPClient().
			Mock("GET", "/api/v4/user",
				func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(`
						{
						  "id": 1,
						  "username": "john_smith",
						  "name": "John Smith",
						  "state": "active",
						  "avatar_url": "http://localhost:3000/uploads/user/avatar/1/cd8.jpeg",
						  "web_url": "http://localhost:3000/john_smith",
						  "created_at": "2012-05-23T08:00:58Z",
						  "bio": "",
						  "bot": false,
						  "location": null,
						  "public_email": "john@example.com",
						  "skype": "",
						  "linkedin": "",
						  "twitter": "",
						  "website_url": "",
						  "organization": "",
						  "job_title": "Operations Specialist",
						  "followers": 1,
						  "following": 1
						}`,
					))
				},
			).
			Mock("GET", "/api/v4/groups/greatuser/projects",
				func(w http.ResponseWriter, r *http.Request) {
					_, _ = w.Write([]byte(`
						[
						  {
							"id": 9,
							"description": "my dotfiles",
							"default_branch": "master",
							"visibility": "internal",
							"ssh_url_to_repo": "git@gitlab.com:greatuser/dotfiles.git",
							"http_url_to_repo": "https://gitlab.com/greatuser/dotfiles",
							"web_url": "https://gitlab.com/greatuser/dotfiles",
							"path": "dotfiles"
						  },
						  {
							"id": 10,
							"description": "my personal website",
							"default_branch": "master",
							"visibility": "internal",
							"ssh_url_to_repo": "git@gitlab.com:greatuser/personal-website.git",
							"http_url_to_repo": "https://gitlab.com/greatuser/personal-website",
							"web_url": "https://gitlab.com/greatuser/personal-website",
							"name": "personal-website",
							"path": "personal-website"
						  }
						]`,
					))
				},
			).Client()
		gitlabClient, err = gitlab.NewClient("sometoken", gitlab.WithHTTPClient(httpClient))
		Expect(err).To(BeNil())
		client = upstream.NewGitlabClient(gitlabClient)
		repositories, err = client.GetRepositories(context.Background(), []string{"greatuser"}, false)
		Expect(err).To(BeNil())
	})
	It("Returns the matching repositories", func() {
		Expect(len(repositories)).To(Equal(2))
		Expect(repositories[0].GetOwner()).To(Equal("greatuser"))
		Expect(repositories[0].GetName()).To(Equal("dotfiles"))
		Expect(repositories[0].GetDescription()).To(Equal("my dotfiles"))
		Expect(repositories[0].GetBrowserHomepageURL()).To(Equal("https://gitlab.com/greatuser/dotfiles"))
		Expect(repositories[0].GetBrowserPullRequestsURL()).To(Equal("https://gitlab.com/greatuser/dotfiles/merge_requests"))
		Expect(repositories[0].GetHTTPSCloneURL()).To(Equal("https://gitlab.com/greatuser/dotfiles"))
		Expect(repositories[0].GetSSHCloneURL()).To(Equal("git@gitlab.com:greatuser/dotfiles.git"))
		Expect(repositories[1].GetName()).To(Equal("personal-website"))
		Expect(repositories[1].GetName()).To(Equal("personal-website"))
		Expect(repositories[1].GetDescription()).To(Equal("my personal website"))
		Expect(repositories[1].GetBrowserHomepageURL()).To(Equal("https://gitlab.com/greatuser/personal-website"))
		Expect(repositories[1].GetBrowserPullRequestsURL()).To(Equal("https://gitlab.com/greatuser/personal-website/merge_requests"))
		Expect(repositories[1].GetHTTPSCloneURL()).To(Equal("https://gitlab.com/greatuser/personal-website"))
		Expect(repositories[1].GetSSHCloneURL()).To(Equal("git@gitlab.com:greatuser/personal-website.git"))
		Expect(repositories[0].GetProvider()).To(Equal("gitlab"))
		Expect(repositories[1].GetProvider()).To(Equal("gitlab"))
	})
})
