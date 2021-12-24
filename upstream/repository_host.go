package upstream

import "context"

type RepositoryHostClient interface {
	GetRepositories(ctx context.Context, owners []string) ([]HostRepository, error)
}

type HostRepository interface {
	GetName() string
	GetOwner() string
	GetDescription() string
	GetBrowserURL() string
	GetCloneURL() string
}
