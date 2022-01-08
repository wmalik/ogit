package service

import (
	"context"

	"ogit/upstream"
)

type Repository struct {
	Name                   string
	Owner                  string
	Description            string
	BrowserHomepageURL     string
	BrowserPullRequestsURL string
	HTTPSCloneURL          string
	SSHCloneURL            string
}

type Repositories []Repository

type RepositoryService struct {
	client         upstream.RepositoryHostClient
	fetchUserRepos bool
}

func NewRepositoryService(client upstream.RepositoryHostClient, fetchUserRepos bool) *RepositoryService {
	return &RepositoryService{client, fetchUserRepos}
}

func (r *RepositoryService) GetRepositoriesByOwners(ctx context.Context, owners []string) (*Repositories, error) {
	repositories, err := r.client.GetRepositories(ctx, owners, r.fetchUserRepos)
	if err != nil {
		return nil, err
	}

	res := make(Repositories, len(repositories))
	for i, repo := range repositories {
		res[i].Owner = repo.GetOwner()
		res[i].Name = repo.GetName()
		res[i].Description = repo.GetDescription()
		res[i].BrowserHomepageURL = repo.GetBrowserHomepageURL()
		res[i].BrowserPullRequestsURL = repo.GetBrowserPullRequestsURL()
		res[i].HTTPSCloneURL = repo.GetHTTPSCloneURL()
		res[i].SSHCloneURL = repo.GetSSHCloneURL()
	}
	return &res, nil
}

func (r *RepositoryService) GetAPIUsage(ctx context.Context) (*APIUsage, error) {
	githubUsage, err := r.client.GetAPIUsage(ctx)
	if err != nil {
		return nil, err
	}

	return &APIUsage{
		Name:          githubUsage.Name,
		Authenticated: githubUsage.Authenticated,
		User:          githubUsage.User,
		Limit:         githubUsage.Limit,
		Remaining:     githubUsage.Remaining,
		ResetsAt:      githubUsage.ResetsAt,
	}, nil
}
