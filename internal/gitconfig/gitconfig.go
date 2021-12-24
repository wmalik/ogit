package gitconfig

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tcnksm/go-gitconfig"
)

type GitConfig struct {
	orgs         []string
	cloneDirPath string
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

	return &GitConfig{orgs: orgs, cloneDirPath: *cloneDirPath}, nil
}

func (c GitConfig) Orgs() []string {
	return c.orgs
}

func (c GitConfig) CloneDirPath() string {
	return c.cloneDirPath
}

func getOrgs() ([]string, error) {
	orgsRaw, err := gitconfig.Entire("ogit.orgs")
	if err != nil {
		return nil, fmt.Errorf("unable to read ogit.orgs from git config: %s", err)
	}

	if orgsRaw == "" {
		return nil, fmt.Errorf("missing ogit.orgs in git config")
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
		if err.Error() != "the key `ogit.clonedirpath` is not found" {
			return nil, fmt.Errorf("unable to read ogit.clonedirpath from git config: %s", err)
		}
	}

	if cloneDirPath == "" {

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		cloneDirPath = path.Join(homeDir, "ogit")

		if err := os.MkdirAll(cloneDirPath, os.ModeDir); err != nil {
			if !strings.HasSuffix(err.Error(), "file exists") {
				return nil, err
			}
		}
	}
	return &cloneDirPath, nil
}
