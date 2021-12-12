package upstream

import (
	"context"
	"log"
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

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(client *github.Client) *GithubClient {
	return &GithubClient{client}
}

func (c *GithubClient) GetRepositories(ctx context.Context, owners []string) []HostRepository {
	opt := &github.RepositoryListOptions{}
	res := []HostRepository{}
	var lock = sync.RWMutex{}
	var wg sync.WaitGroup
	wg.Add(len(owners))
	for _, owner := range owners {
		go func(owner string) {
			defer wg.Done()
			repos, _, err := c.client.Repositories.List(ctx, owner, opt)
			if err != nil {
				log.Printf("error while fetching repositories on Github: %s", err)
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
	return res
}
