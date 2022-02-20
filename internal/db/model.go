package db

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	Title                  string `gorm:"unique"`
	Owner                  string
	Name                   string
	Description            string
	BrowserHomepageURL     string
	BrowserPullRequestsURL string
	HTTPSCloneURL          string
	SSHCloneURL            string
}

func NewRepository(
	title,
	owner,
	name,
	description,
	browserHomepageURL,
	browserPullRequestsURL,
	httpsCloneURL,
	sshCloneURL string,
) Repository {
	return Repository{
		Title:                  title,
		Owner:                  owner,
		Name:                   name,
		Description:            description,
		BrowserHomepageURL:     browserHomepageURL,
		BrowserPullRequestsURL: browserPullRequestsURL,
		HTTPSCloneURL:          httpsCloneURL,
		SSHCloneURL:            sshCloneURL,
	}
}
