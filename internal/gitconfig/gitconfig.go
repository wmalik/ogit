package gitconfig

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tcnksm/go-gitconfig"
)

type GitConfig struct {
	orgs         []string
	cloneDirPath string
	// whether to fetch repos associated with the authenticated user
	fetchAuthenticatedUserRepos bool
	useSSHAgent                 bool
	// the path to the SSH private key used for git operations e.g. clone
	privKeyPath string
}

// readGitConfig loads the value of ogit.orgs from ~/.gitconfig
func ReadGitConfig() (*GitConfig, error) {

	orgs, err := getOrgs()
	if err != nil {
		return nil, err
	}

	cloneDirPath, err := getCloneDirPath()
	if err != nil {
		return nil, err
	}

	fetchUserRepos, err := getFetchAuthenticatedUserRepos()
	if err != nil {
		return nil, err
	}

	useSSHAgent, err := getUseSSHAgent()
	if err != nil {
		return nil, err
	}

	privKeyPath, err := getPrivKeyPath()
	if err != nil {
		return nil, err
	}

	return &GitConfig{
		orgs:                        orgs,
		cloneDirPath:                *cloneDirPath,
		fetchAuthenticatedUserRepos: fetchUserRepos,
		useSSHAgent:                 useSSHAgent,
		privKeyPath:                 privKeyPath,
	}, nil
}

func (c GitConfig) Orgs() []string {
	return c.orgs
}

func (c GitConfig) CloneDirPath() string {
	return c.cloneDirPath
}

func (c GitConfig) FetchAuthenticatedUserRepos() bool {
	return c.fetchAuthenticatedUserRepos
}

func (c GitConfig) UseSSHAgent() bool {
	return c.useSSHAgent
}

func (c GitConfig) PrivKeyPath() string {
	return c.privKeyPath
}

func getOrgs() ([]string, error) {
	orgsRaw, err := gitconfig.Entire("ogit.orgs")
	if err != nil {
		return nil, fmt.Errorf("missing manadatory config in git config: %s", err)
	}

	if orgsRaw == "" {
		return nil, errors.New("blank ogit.orgs in git config")
	}

	orgs := []string{}
	for _, org := range strings.Split(orgsRaw, ",") {
		orgs = append(orgs, strings.TrimSpace(org))
	}

	return orgs, err
}

func getCloneDirPath() (*string, error) {
	var cloneDirPath string
	var err error
	cloneDirPath, err = gitconfig.Entire("ogit.clonedirpath")
	if err != nil {
		return nil, fmt.Errorf("missing ogit.clonedirpath in git config: %s", err)
	}

	return &cloneDirPath, nil
}

func getFetchAuthenticatedUserRepos() (bool, error) {
	fetchAuthenticatedUserRepos, err := gitconfig.Entire("ogit.fetchAuthenticatedUserRepos")
	if err != nil {
		return false, fmt.Errorf("missing ogit.fetchAuthenticatedUserRepos in git config: %s", err)
	}

	if strings.TrimSpace(fetchAuthenticatedUserRepos) == "false" {
		return false, nil
	}

	return true, nil
}

func getUseSSHAgent() (bool, error) {
	useSSHAgent, err := gitconfig.Entire("ogit.useSSHAgent")
	if err != nil {
		return false, fmt.Errorf("missing ogit.useSSHAgent in git config: %s", err)
	}

	if strings.TrimSpace(useSSHAgent) == "false" {
		return false, nil
	}

	return true, nil
}

func getPrivKeyPath() (string, error) {
	privKeyPath, err := gitconfig.Entire("ogit.privKeyPath")
	if err != nil {
		return "", fmt.Errorf("missing ogit.privKeyPath in git config: %s", err)
	}

	return strings.TrimSpace(privKeyPath), nil
}
