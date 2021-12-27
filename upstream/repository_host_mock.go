package upstream

import (
	"context"
)

type MockRepository struct {
	Owner       string
	Name        string
	Description string
	BrowserURL  string
	CloneURL    string
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

func (r *MockRepository) GetBrowserURL() string {
	return r.BrowserURL
}

func (r *MockRepository) GetCloneURL() string {
	return r.CloneURL
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

func (c *MockClient) GetRepositories(ctx context.Context, owners []string) []HostRepository {
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
	return res
}
