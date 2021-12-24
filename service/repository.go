package service

import (
	"context"

	"ogit/upstream"
)

type Repository struct {
	Name        string
	Owner       string
	Description string
	BrowserURL  string
	CloneURL    string
}

type Repositories []Repository

type RepositoryService struct {
	client upstream.RepositoryHostClient
}

func NewRepositoryService(client upstream.RepositoryHostClient) *RepositoryService {
	return &RepositoryService{client}
}

func (r *RepositoryService) GetRepositoriesByOwners(ctx context.Context, owners []string) (*Repositories, error) {
	repositories, err := r.client.GetRepositories(ctx, owners)
	if err != nil {
		return nil, err
	}

	res := make(Repositories, len(repositories))
	for i, repo := range repositories {
		res[i].Owner = repo.GetOwner()
		res[i].Name = repo.GetName()
		res[i].Description = repo.GetDescription()
		res[i].BrowserURL = repo.GetBrowserURL()
		res[i].CloneURL = repo.GetCloneURL()
	}
	return &res, nil
}
