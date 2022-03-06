package sync

import (
	"context"
	"fmt"
	"log"
	"ogit/internal/db"
	"ogit/internal/gitconfig"
	"ogit/service"
	"ogit/upstream"
	"os"
	"path"
)

// Sync fetches the repository metadata from upstream and stores it in the local
// database (on disk)
func Sync(ctx context.Context, gitConf *gitconfig.GitConfig) error {
	gitlabClient, err := upstream.NewGitlabClientWithToken(os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		log.Fatalln(err)
	}

	rs := service.NewRepositoryService(
		upstream.NewGithubClientWithToken(os.Getenv("GITHUB_TOKEN")),
		gitlabClient,
		gitConf.FetchAuthenticatedUserRepos(),
	)

	log.Println("Syncing repositories")
	repos, err := rs.GetRepositoriesByOwners(ctx, gitConf.Orgs(), gitConf.GitlabGroups())
	if err != nil {
		log.Fatalln(err)
	}

	localDB, err := db.NewDB(path.Join(gitConf.StoragePath(), "ogit.db"))
	if err != nil {
		log.Fatalln(err)
	}

	if err := localDB.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := localDB.UpsertRepositories(ctx, toDatabaseRepositories(repos)); err != nil {
		log.Fatalln(err)
	}

	return nil
}

func toDatabaseRepositories(repos *service.Repositories) []db.Repository {
	dbRepos := []db.Repository{}
	for _, repo := range *repos {
		dbRepos = append(dbRepos, db.NewRepository(
			fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
			repo.Owner,
			repo.Name,
			repo.Description,
			repo.BrowserHomepageURL,
			repo.BrowserPullRequestsURL,
			repo.HTTPSCloneURL,
			repo.SSHCloneURL,
		),
		)
	}

	return dbRepos

}
