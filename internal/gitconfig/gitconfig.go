package gitconfig

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tcnksm/go-gitconfig"
)

type GitConfig struct {
	orgs         []string
	gitlabGroups []string
	storagePath  string
	// whether to fetch repos associated with the authenticated user
	fetchUserRepos bool
	useSSHAgent    bool
	// the path to the SSH private key used for git operations e.g. clone
	privKeyPath string
}

func defaultGitConfig() *GitConfig {
	return &GitConfig{
		storagePath:    filepath.Join(os.TempDir(), "ogit"),
		fetchUserRepos: true,
		useSSHAgent:    true,
		privKeyPath:    "",
	}
}

// readGitConfig loads the value of ogit.orgs from ~/.gitconfig
func ReadGitConfig() (*GitConfig, error) {

	conf := defaultGitConfig()

	orgs, err := getOrgs()
	if err != nil {
		return nil, err
	}
	if len(orgs) > 0 {
		conf.orgs = orgs
	}

	gitlabGroups, err := getGitlabGroups()
	if err != nil {
		return nil, err
	}
	if len(gitlabGroups) > 0 {
		conf.gitlabGroups = gitlabGroups
	}

	storagePath, err := getStoragePath()
	if err != nil {
		return nil, err
	}
	if storagePath != "" {
		conf.storagePath = storagePath
	}

	fetchUserRepos, err := getFetchAuthenticatedUserRepos()
	if err != nil {
		return nil, err
	}
	if fetchUserRepos != true {
		conf.fetchUserRepos = fetchUserRepos
	}

	useSSHAgent, privKeyPath, err := getSSHAuth()
	if err != nil {
		return nil, err
	}
	if useSSHAgent != true {
		conf.useSSHAgent = useSSHAgent
	}
	if privKeyPath != "" {
		conf.privKeyPath = privKeyPath
	}

	return conf, nil
}

func (c GitConfig) Orgs() []string {
	return c.orgs
}

func (c GitConfig) GitlabGroups() []string {
	return c.gitlabGroups
}

func (c GitConfig) StoragePath() string {
	return c.storagePath
}

func (c GitConfig) FetchUserRepos() bool {
	return c.fetchUserRepos
}

func (c GitConfig) UseSSHAgent() bool {
	return c.useSSHAgent
}

func (c GitConfig) PrivKeyPath() string {
	return c.privKeyPath
}

func getOrgs() ([]string, error) {
	orgsRaw, err := gitconfig.Entire("ogit.github.orgs")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return []string{}, nil
		}
		return nil, err
	}

	orgs := []string{}
	for _, org := range strings.Split(orgsRaw, ",") {
		if org != "" {
			orgs = append(orgs, strings.TrimSpace(org))
		}
	}

	return orgs, err
}

func getGitlabGroups() ([]string, error) {
	gitlabGroupsRaw, err := gitconfig.Entire("ogit.gitlab.orgs")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return []string{}, nil
		}
		return nil, err
	}

	gitlabGroups := []string{}
	for _, org := range strings.Split(gitlabGroupsRaw, ",") {
		if org != "" {
			gitlabGroups = append(gitlabGroups, strings.TrimSpace(org))
		}
	}

	return gitlabGroups, err
}

func getStoragePath() (string, error) {
	var storagePath string
	var err error
	storagePath, err = gitconfig.Entire("ogit.storagePath")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return "", nil
		}
		return "", err
	}

	return storagePath, nil
}

func getFetchAuthenticatedUserRepos() (bool, error) {
	fetchUserRepos, err := gitconfig.Entire("ogit.fetchUserRepos")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return true, nil
		}
		return false, err
	}

	if strings.TrimSpace(fetchUserRepos) == "false" {
		return false, nil
	}

	return true, nil
}

func getSSHAuth() (useSSHAgent bool, privKeyPath string, err error) {
	sshAuth, err := gitconfig.Entire("ogit.sshAuth")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return true, "", nil
		}
		return false, "", err
	}

	if strings.TrimSpace(sshAuth) == "ssh-agent" {
		return true, "", nil
	}

	return false, strings.TrimSpace(sshAuth), nil
}
