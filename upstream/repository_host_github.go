package upstream

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

type GithubRepository struct {
	github.Repository
}

func (r *GithubRepository) GetName() string {
	if r.Name == nil {
		return ""
	}
	return *r.Name
}

func (r *GithubRepository) GetOwner() string {
	if r.Owner == nil {
		return ""
	}
	return r.Owner.GetLogin()
}

func (r *GithubRepository) GetDescription() string {
	if r.Description == nil {
		return ""
	}
	return *r.Description
}

func (r *GithubRepository) GetBrowserURL() string {
	if r.GetHTMLURL() == "" {
		return ""
	}

	return r.GetHTMLURL()
}

func (r *GithubRepository) GetCloneURL() string {
	if r.Repository.GetCloneURL() == "" {
		return ""
	}

	return r.Repository.GetCloneURL()
}

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(client *github.Client) *GithubClient {
	return &GithubClient{client}
}

func (c *GithubClient) GetRepositories(ctx context.Context, owners []string) ([]HostRepository, error) {
	opt := &github.RepositoryListOptions{Sort: "updated"}
	res := []HostRepository{}
	var lock = sync.RWMutex{}
	var wg sync.WaitGroup
	var errResult error
	wg.Add(len(owners))
	for _, owner := range owners {
		go func(owner string) {
			defer wg.Done()
			repos, _, err := c.client.Repositories.List(ctx, owner, opt)
			if err != nil {
				errResult = fmt.Errorf("error while fetching repositories on Github: %s", err)
				return
			}
			grepos := make([]HostRepository, len(repos))
			for i, r := range repos {
				grepos[i] = &GithubRepository{*r}
			}
			lock.Lock()
			res = append(res, grepos...)
			lock.Unlock()
		}(owner)
	}
	wg.Wait()
	return res, errResult
}
