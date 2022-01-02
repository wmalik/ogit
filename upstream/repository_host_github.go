package upstream

import (
	"context"
	"fmt"
	"math"
	"log"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const pageSize = 100

type GithubRepository struct {
	github.Repository
}

type APIUsage struct {
	Name          string
	Authenticated bool
	User          string
	Limit         int
	Remaining     int
	ResetsAt      time.Time
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

func NewGithubClientWithToken(token string) *GithubClient {
	if token == "" {
		return &GithubClient{github.NewClient(nil)}
	}

	return &GithubClient{
		github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
			),
		),
	}
}

func (c *GithubClient) GetRepositories(ctx context.Context, owners []string) ([]HostRepository, error) {
	res := []HostRepository{}
	var m sync.Map

	var g errgroup.Group

	for _, owner := range owners {
		g.Go(func(owner string) func() error {
			return func() error {
				repos, err := c.getRepositoriesForOwner(ctx, owner, 0)
				if err != nil {
					return err
				}

				m.Store(owner, repos)
				return nil
			}
		}(owner))
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	m.Range(func(key, value interface{}) bool {
		res = append(res, value.([]HostRepository)...)
		return true
	})

	return res, nil
}

func (c *GithubClient) GetAPIUsage(ctx context.Context) (*APIUsage, error) {
	limits, _, err := c.client.RateLimits(ctx)
	if err != nil {
		return nil, err
	}

	usage := &APIUsage{
		Name:      "GitHub",
		Limit:     limits.GetCore().Limit,
		Remaining: limits.GetCore().Remaining,
		ResetsAt:  limits.GetCore().Reset.Time,
	}

	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		log.Println("Unable to get user information, perhaps a github token is not set?")
		usage.Authenticated = false
	} else {
		usage.Authenticated = true
		usage.User = user.GetLogin()
	}

	return usage, nil
}

func (c *GithubClient) getRepositoriesForOwner(ctx context.Context, owner string, startPage int) ([]HostRepository, error) {
	var reposAcc []*github.Repository
	opt := &github.RepositoryListOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			Page:    startPage,
			PerPage: pageSize,
		},
	}

	for {
		repos, resp, err := c.client.Repositories.List(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		reposAcc = append(reposAcc, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	repos := make([]HostRepository, len(reposAcc))
	for i, r := range reposAcc {
		repos[i] = &GithubRepository{*r}
	}
	return repos, nil
}

func (c *GithubClient) GetRateLimits(ctx context.Context) (string, error) {
	limits, _, err := c.client.RateLimits(ctx)
	if err != nil {
		return "", fmt.Errorf("error while fetching github rate limits: %s", err)
	}
	return fmt.Sprintf("[GitHub API Usage (%d of %d) (resets in %d mins)]",
		(limits.GetCore().Limit - limits.GetCore().Remaining),
		limits.GetCore().Limit,
		int(math.Ceil(time.Until(limits.GetCore().Reset.Time).Minutes())),
	), nil
}
