package upstream

import (
	"context"
)

type MockRepository struct {
	Provider               string
	Owner                  string
	Name                   string
	Description            string
	BrowserHomepageURL     string
	BrowserPullRequestsURL string
	HTTPSCloneURL          string
	SSHCloneURL            string
}

func (r *MockRepository) GetProvider() string {
	return r.Provider
}

func (r *MockRepository) GetName() string {
	return r.Name
}

func (r *MockRepository) GetOwner() string {
	return r.Owner
}

func (r *MockRepository) GetDescription() string {
	return r.Description
}

func (r *MockRepository) GetBrowserHomepageURL() string {
	return r.BrowserHomepageURL
}

func (r *MockRepository) GetBrowserPullRequestsURL() string {
	return r.BrowserPullRequestsURL
}

func (r *MockRepository) GetHTTPSCloneURL() string {
	return r.HTTPSCloneURL
}

func (r *MockRepository) GetSSHCloneURL() string {
	return r.SSHCloneURL
}

type MockClient struct {
	repositories []MockRepository
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (c *MockClient) WithRepositories(repos []MockRepository) *MockClient {
	c.repositories = repos
	return c
}

func (c *MockClient) GetRepositories(ctx context.Context, owners []string, fetchAuthenticationUserRepos bool) ([]HostRepository, error) {
	inputOwners := map[string]struct{}{}
	for _, owner := range owners {
		inputOwners[owner] = struct{}{}
	}
	res := []HostRepository{}
	for i, repo := range c.repositories {
		if _, found := inputOwners[repo.GetOwner()]; found {
			res = append(res, &c.repositories[i])
		}
	}
	return res, nil
}
