package upstream

import (
	"context"
)

type RepositoryHostClient interface {
	GetRepositories(ctx context.Context, owners []string, fetchAuthenticatedUserRepos bool) ([]HostRepository, error)
	GetAPIUsage(ctx context.Context) (*APIUsage, error)
}

type HostRepository interface {
	GetName() string
	GetOwner() string
	GetDescription() string
	GetBrowserHomepageURL() string
	GetBrowserPullRequestsURL() string
	GetHTTPSCloneURL() string
	GetSSHCloneURL() string
}
