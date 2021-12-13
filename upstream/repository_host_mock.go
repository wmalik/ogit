package upstream

import (
	"context"
)

type MockRepository struct {
	Owner string
	Name  string
}

func (r *MockRepository) GetName() string {
	return r.Name
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
		if _, found := inputOwners[repo.Owner]; found {
			res = append(res, &c.repositories[i])
		}
	}
	return res
}
