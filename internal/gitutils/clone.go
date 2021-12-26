package gitutils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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

func (r *Repository) String() string {
	return fmt.Sprintf("%s -> %s (%s %s)", r.GitURL, r.Path, r.HeadRef[:6], r.HeadRefName)
}

func (r *Repository) LastCommit() string {
	return fmt.Sprintf(
		"%s %s %s (%s) %s (%s)",
		r.HeadRef,
		r.HeadRefName,
		r.LastCommitInfo.AuthorName,
		r.LastCommitInfo.AuthorEmail,
		r.LastCommitInfo.Message,
		r.LastCommitInfo.When,
	)
}

func CloneToDisk(ctx context.Context, gitURL, path string, progress io.Writer) (*Repository, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}
	repo, err := git.PlainCloneContext(ctx, path, false,
		&git.CloneOptions{
			URL:      gitURL,
			Progress: progress,
			Depth:    1,
		},
	)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil, ErrRepoAlreadyCloned
		}
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commitObject, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	return &Repository{
		GitURL:      gitURL,
		Path:        path,
		HeadRefName: head.Name().Short(),
		HeadRef:     head.Hash().String(),
		LastCommitInfo: commitInfo{
			Message:     commitObject.Message,
			AuthorName:  commitObject.Author.Name,
			AuthorEmail: commitObject.Author.Email,
			When:        commitObject.Author.When,
		},
	}, nil
}

// Cloned checks if a path contains a .git directory
func Cloned(dir string) bool {
	if _, err := os.Stat(path.Join(dir, ".git")); err != nil {
		return false
	}

	return true
}
