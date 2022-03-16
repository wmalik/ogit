package upstream

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

const pageSize = 100
const githubUpstream = "github.com"

type GithubRepository struct {
	github.Repository
}

func (r *GithubRepository) GetProvider() string {
	return "github"
}

func (r *GithubRepository) GetName() string {
	return r.Repository.GetName()
}

func (r *GithubRepository) GetOwner() string {
	if r.Repository.GetOwner() == nil {
		return ""
	}
	return r.Repository.GetOwner().GetLogin()
}

func (r *GithubRepository) GetDescription() string {
	return r.Repository.GetDescription()
}

func (r *GithubRepository) GetBrowserHomepageURL() string {
	return r.GetHTMLURL()
}

func (r *GithubRepository) GetBrowserPullRequestsURL() string {
	return r.GetHTMLURL() + "/pulls"
}

func (r *GithubRepository) GetOrgURL() string {
	parsed, err := url.Parse(r.GetHTMLURL())
	if err != nil {
		log.Println("unable to parse org url")
		return ""
	}
	parsed.Path = filepath.Dir(parsed.Path)
	return parsed.String()
}

func (r *GithubRepository) GetIssuesURL() string {
	return r.GetHTMLURL() + "/issues"
}

func (r *GithubRepository) GetCIURL() string {
	return r.GetHTMLURL() + "/actions"
}

func (r *GithubRepository) GetReleasesURL() string {
	return r.GetHTMLURL() + "/releases"
}

func (r *GithubRepository) GetSettingsURL() string {
	return r.GetHTMLURL() + "/settings"
}

func (r *GithubRepository) GetHTTPSCloneURL() string {
	return r.Repository.GetHTMLURL()
}

func (r *GithubRepository) GetSSHCloneURL() string {
	return r.Repository.GetSSHURL()
}

type GithubClient struct {
	client   *github.Client
	username string
}

func NewGithubClient(client *github.Client) *GithubClient {
	return &GithubClient{client: client, username: "nobody"}
}

func NewGithubClientWithToken(token string) *GithubClient {
	if token == "" {
		return &GithubClient{client: github.NewClient(nil), username: "nobody"}
	}

	return &GithubClient{
		client: github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
			),
		),
		username: "nobody",
	}
}

func (c *GithubClient) GetRepositories(ctx context.Context, owners []string, fetchUserRepos bool) ([]HostRepository, error) {
	res := HostRepositories{}
	var m sync.Map

	if fetchUserRepos {
		owners = append(owners, "")
	}

	if err := c.setUserInfo(ctx); err != nil {
		return nil, err
	}

	logAuthenticatedUser(githubUpstream, c.username)

	var g errgroup.Group

	for _, owner := range owners {
		g.Go(func(owner string) func() error {
			return func() error {
				var repos []HostRepository
				var err error
				repos, err = c.getRepositoriesForOwner(ctx, owner, 0)
				if err != nil {
					if err.Error() != "not found" {
						return err
					}

				}

				if len(repos) == 0 {
					repos, err = c.getRepositoriesForOrg(ctx, owner, 0)
					if err != nil {
						return err
					}
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

	return res.DeDuplicate(), nil
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

		logPaginationStatus(githubUpstream, owner, len(repos), resp.LastPage-resp.NextPage, strconv.Itoa(resp.Remaining))

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

func (c *GithubClient) getRepositoriesForOrg(ctx context.Context, org string, startPage int) ([]HostRepository, error) {
	var reposAcc []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			Page:    startPage,
			PerPage: pageSize,
		},
	}

	for {
		repos, resp, err := c.client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return nil, err
			}
			return []HostRepository{}, nil
		}

		logPaginationStatus(githubUpstream, org, len(repos), resp.LastPage-resp.NextPage, strconv.Itoa(resp.Remaining))

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

// setUserInfo fetches the authenticated user's information and stores it.
func (c *GithubClient) setUserInfo(ctx context.Context) error {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		log.Println("Unable to get user information, perhaps a github token is not set?")
		return err
	}

	c.username = user.GetLogin()
	return nil
}
