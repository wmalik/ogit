package service

import (
	"context"

	"ogit/upstream"
)

type Repository struct {
	Name string
}

type Repositories []Repository

type RepositoryService struct {
	client upstream.RepositoryHostClient
}

func NewRepositoryService(client upstream.RepositoryHostClient) *RepositoryService {
	return &RepositoryService{client}
}

func (r *RepositoryService) GetRepositoriesByOwners(ctx context.Context, owners []string) Repositories {
	repositories := r.client.GetRepositories(ctx, owners)
	res := make(Repositories, len(repositories))
	for i, repo := range repositories {
		res[i].Name = repo.GetName()
	}
	return res
}
