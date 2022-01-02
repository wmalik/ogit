package upstream

import "context"

type RepositoryHostClient interface {
	GetRepositories(ctx context.Context, owners []string) ([]HostRepository, error)
	GetRateLimits(ctx context.Context) (string, error)
	GetAPIUsage(ctx context.Context) (*APIUsage, error)
}

type HostRepository interface {
	GetName() string
	GetOwner() string
	GetDescription() string
	GetBrowserURL() string
	GetCloneURL() string
}
