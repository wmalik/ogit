package gitutils

import (
	"strings"

	"github.com/go-git/go-git/v5"
)

func ReadRepository(path string) (*Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
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
		Path:        path,
		HeadRefName: head.Name().Short(),
		HeadRef:     head.Hash().String(),
		LastCommitInfo: commitInfo{
			Message:     strings.TrimSpace(strings.ReplaceAll(commitObject.Message, "\n", " ")),
			AuthorName:  commitObject.Author.Name,
			AuthorEmail: commitObject.Author.Email,
			When:        commitObject.Author.When,
		},
	}, nil

}
