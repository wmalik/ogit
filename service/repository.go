package service

import (
	"context"

	"github.com/wmalik/ogit/upstream"
)

type Repository struct {
	Provider               string
	Name                   string
	Owner                  string
	Description            string
	BrowserHomepageURL     string
	BrowserPullRequestsURL string
	HTTPSCloneURL          string
	SSHCloneURL            string
	OrgURL                 string
	IssuesURL              string
	CIURL                  string
	ReleasesURL            string
	SettingsURL            string
}

type Repositories []Repository

type RepositoryService struct {
	client         upstream.RepositoryHostClient
	gitlabClient   upstream.RepositoryHostClient
	fetchUserRepos bool
}

func NewRepositoryService(client upstream.RepositoryHostClient, gitlabClient upstream.RepositoryHostClient, fetchUserRepos bool) *RepositoryService {
	return &RepositoryService{client, gitlabClient, fetchUserRepos}
}

func (r *RepositoryService) GetRepositoriesByOwners(ctx context.Context, owners []string, gitlabOwners []string) (*Repositories, error) {
	repositories, err := r.client.GetRepositories(ctx, owners, r.fetchUserRepos)
	if err != nil {
		return nil, err
	}

	gitlabRepositories, err := r.gitlabClient.GetRepositories(ctx, gitlabOwners, r.fetchUserRepos)
	if err != nil {
		return nil, err
	}

	allRepositories := append(repositories, gitlabRepositories...)

	res := make(Repositories, len(allRepositories))
	for i, repo := range allRepositories {
		res[i].Provider = repo.GetProvider()
		res[i].Owner = repo.GetOwner()
		res[i].Name = repo.GetName()
		res[i].Description = repo.GetDescription()
		res[i].BrowserHomepageURL = repo.GetBrowserHomepageURL()
		res[i].BrowserPullRequestsURL = repo.GetBrowserPullRequestsURL()
		res[i].OrgURL = repo.GetOrgURL()
		res[i].IssuesURL = repo.GetIssuesURL()
		res[i].CIURL = repo.GetCIURL()
		res[i].ReleasesURL = repo.GetReleasesURL()
		res[i].SettingsURL = repo.GetSettingsURL()
		res[i].HTTPSCloneURL = repo.GetHTTPSCloneURL()
		res[i].SSHCloneURL = repo.GetSSHCloneURL()
	}
	return &res, nil
}
