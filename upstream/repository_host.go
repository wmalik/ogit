package upstream

import (
	"context"
	"fmt"
)

type RepositoryHostClient interface {
	GetRepositories(ctx context.Context, owners []string, fetchUserRepos bool) ([]HostRepository, error)
}

type HostRepository interface {
	GetProvider() string
	GetName() string
	GetOwner() string
	GetDescription() string
	GetBrowserHomepageURL() string
	GetBrowserPullRequestsURL() string
	GetHTTPSCloneURL() string
	GetSSHCloneURL() string
	GetOrgURL() string
	GetIssuesURL() string
	GetCIURL() string
	GetReleasesURL() string
	GetSettingsURL() string
}

type HostRepositories []HostRepository

func (hr HostRepositories) DeDuplicate() []HostRepository {
	var results []HostRepository
	uniqueMap := map[string]HostRepository{}
	for _, hostRepo := range hr {
		key := fmt.Sprintf("%s/%s", hostRepo.GetOwner(), hostRepo.GetName())
		if _, ok := uniqueMap[key]; !ok {
			uniqueMap[key] = hostRepo
			results = append(results, hostRepo)
		}
	}

	return results
}
