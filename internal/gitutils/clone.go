package gitutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var ErrRepoAlreadyCloned error = errors.New("Repository already cloned")

type commitInfo struct {
	AuthorName  string
	AuthorEmail string
	Message     string
	When        time.Time
}

type Repository struct {
	GitURL         string
	Path           string
	HeadRefName    string
	HeadRef        string
	LastCommitInfo commitInfo
}

type GitUtils struct {
	auth           ssh.AuthMethod
	cloneOverHTTPS bool
}

func NewGitUtils(useSSHAgent bool, privKeyPath string) (*GitUtils, error) {
	if privKeyPath != "" {
		return newGitUtilsWithPrivKey(privKeyPath)
	}

	if useSSHAgent {
		return newGitUtilsWithSSHAgent()
	}

	return &GitUtils{auth: nil, cloneOverHTTPS: true}, nil
}

func newGitUtilsWithSSHAgent() (*GitUtils, error) {
	auth, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		return nil, err
	}
	return &GitUtils{auth: auth, cloneOverHTTPS: false}, nil
}

func newGitUtilsWithPrivKey(privKeyPath string) (*GitUtils, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", privKeyPath, "")
	if err != nil {
		return nil, err
	}

	return &GitUtils{auth: auth, cloneOverHTTPS: false}, nil
}

func (r *Repository) String() string {
	return fmt.Sprintf("%s -> %s (%s %s)", r.GitURL, r.Path, r.HeadRef[:6], r.HeadRefName)
}

func (r *Repository) LastCommit() string {
	return fmt.Sprintf(
		"%s %s %s (%s) %s (%s)",
		r.HeadRef[:6],
		r.HeadRefName,
		r.LastCommitInfo.Message,
		r.LastCommitInfo.AuthorName,
		r.LastCommitInfo.AuthorEmail,
		r.LastCommitInfo.When.Format("January 2, 2006"),
	)
}

// CloneToDisk clones a repository to a path on disk.
// The repository is cloned first to a temporary path, and then renamed to the
// desired path. The function guarantees that if `path` exists, it contains
// a fully cloned repository.
// If an authentication method has been configured, the repository is cloned
// using sshURL, otherwise it is cloned using httpsURL.  The progress of the
// clone operation is streamed to the progress io.Writer
func (gu *GitUtils) CloneToDisk(ctx context.Context, httpsURL, sshURL, path string, progress io.Writer) (string, error) {
	cloneURL := sshURL
	if gu.cloneOverHTTPS {
		cloneURL = httpsURL
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	repo, err := git.PlainCloneContext(ctx, tmpDir, false,
		&git.CloneOptions{
			URL:      cloneURL,
			Progress: progress,
			Depth:    1,
			Auth:     gu.auth,
		},
	)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	commitObject, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", err
	}

	repository := &Repository{
		GitURL:      cloneURL,
		Path:        path,
		HeadRefName: head.Name().Short(),
		HeadRef:     head.Hash().String(),
		LastCommitInfo: commitInfo{
			Message:     strings.TrimSpace(commitObject.Message),
			AuthorName:  commitObject.Author.Name,
			AuthorEmail: commitObject.Author.Email,
			When:        commitObject.Author.When,
		},
	}

	if err := os.Rename(tmpDir, path); err != nil {
		return "", fmt.Errorf("rename failed after cloning: %s", err)
	}

	return repository.String(), nil
}

// Cloned checks if a path contains a .git directory
func Cloned(dir string) bool {
	if _, err := os.Stat(path.Join(dir, ".git")); err != nil {
		return false
	}

	return true
}
