package upstream

import (
	"context"
	"log"
	"sync"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

const gitlabPageSize = 100
const gitlabUpstream = "gitlab.com"

type GitlabProject struct {
	gitlab.Project
	Username string
}

func (r *GitlabProject) GetName() string {
	return r.Project.Path
}

func (r *GitlabProject) GetOwner() string {
	return r.Username
}

func (r *GitlabProject) GetDescription() string {
	return r.Project.Description
}

func (r *GitlabProject) GetBrowserHomepageURL() string {
	return r.Project.WebURL
}

func (r *GitlabProject) GetBrowserPullRequestsURL() string {
	return r.Project.HTTPURLToRepo + "/merge_requests"
}

func (r *GitlabProject) GetHTTPSCloneURL() string {
	return r.Project.HTTPURLToRepo
}

func (r *GitlabProject) GetSSHCloneURL() string {
	return r.Project.SSHURLToRepo
}

type GitlabClient struct {
	client   *gitlab.Client
	username string // the username of authenticated user
	userID   int    // the id of authenticated user
}

func NewGitlabClient(client *gitlab.Client) *GitlabClient {
	return &GitlabClient{client: client, username: "nobody"}
}

func NewGitlabClientWithToken(token string) (*GitlabClient, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &GitlabClient{client: client, username: "nobody"}, nil
}

func (c *GitlabClient) GetRepositories(ctx context.Context, groups []string, fetchUserRepos bool) ([]HostRepository, error) {
	res := HostRepositories{}
	var m sync.Map

	if err := c.setUserInfo(ctx); err != nil {
		return nil, err
	}

	logAuthenticatedUser(gitlabUpstream, c.username)

	var g errgroup.Group
	if fetchUserRepos {
		g.Go(func(ctx context.Context) func() error {
			return func() error {
				userProjects, err := c.getProjectsForAuthUser(ctx, c.userID, c.username)
				if err != nil {
					return err
				}

				m.Store(c.username, userProjects)
				return nil
			}
		}(ctx))
	}

	for _, group := range groups {
		g.Go(func(group string) func() error {
			return func() error {
				repos, err := c.getProjectsForGroup(ctx, group)
				if err != nil {
					return err
				}

				m.Store(group, repos)
				return nil
			}
		}(group))
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

func (c *GitlabClient) getProjectsForAuthUser(ctx context.Context, userID int, username string) ([]HostRepository, error) {
	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: gitlabPageSize,
			Page:    1,
		},
	}

	var allProjects []*gitlab.Project
	for {
		// Get the first page with projects.
		projects, resp, err := c.client.Projects.ListUserProjects(userID, opt)
		if err != nil {
			log.Fatal(err)
		}

		logPaginationStatus(gitlabUpstream, username, len(projects), resp.TotalPages-resp.NextPage-1, resp.Header.Get("RateLimit-Remaining"))

		allProjects = append(allProjects, projects...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	repos := make([]HostRepository, len(allProjects))
	for i, p := range allProjects {
		repos[i] = &GitlabProject{Project: *p, Username: username}
	}
	return repos, nil
}

func (c *GitlabClient) getProjectsForGroup(ctx context.Context, group string) ([]HostRepository, error) {
	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: gitlabPageSize,
			Page:    1,
		},
	}

	var allProjects []*gitlab.Project
	for {
		groupProjects, resp, err := c.client.Groups.ListGroupProjects(
			group,
			opt,
		)
		if err != nil {
			return nil, err
		}

		logPaginationStatus(gitlabUpstream, group, len(groupProjects), resp.TotalPages-resp.NextPage-1, resp.Header.Get("RateLimit-Remaining"))

		allProjects = append(allProjects, groupProjects...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	repos := make([]HostRepository, len(allProjects))
	for i, p := range allProjects {
		repos[i] = &GitlabProject{*p, group}
	}
	return repos, nil
}

// setUserInfo fetches the authenticated user's information and stores it
func (c *GitlabClient) setUserInfo(ctx context.Context) error {
	user, _, err := c.client.Users.CurrentUser()
	if err != nil {
		return err
	}

	c.username = user.Username
	c.userID = user.ID
	return nil
}
