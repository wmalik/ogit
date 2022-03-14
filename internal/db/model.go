package db

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	Provider               string
	Title                  string `gorm:"unique"`
	Owner                  string
	Name                   string
	Description            string
	BrowserHomepageURL     string
	BrowserPullRequestsURL string
	OrgURL                 string
	IssuesURL              string
	CIURL                  string
	ReleasesURL            string
	SettingsURL            string
	HTTPSCloneURL          string
	SSHCloneURL            string
}

func NewRepository(
	provider,
	title,
	owner,
	name,
	description,
	browserHomepageURL,
	browserPullRequestsURL,
	orgURL,
	issuesURL,
	ciURL,
	releasesURL,
	settingsURL,
	httpsCloneURL,
	sshCloneURL string,
) Repository {
	return Repository{
		Provider:               provider,
		Title:                  title,
		Owner:                  owner,
		Name:                   name,
		Description:            description,
		BrowserHomepageURL:     browserHomepageURL,
		BrowserPullRequestsURL: browserPullRequestsURL,
		OrgURL:                 orgURL,
		IssuesURL:              issuesURL,
		CIURL:                  ciURL,
		ReleasesURL:            releasesURL,
		SettingsURL:            settingsURL,
		HTTPSCloneURL:          httpsCloneURL,
		SSHCloneURL:            sshCloneURL,
	}
}
