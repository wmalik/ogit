package browser

import "ogit/service"

type refreshReposMsg struct{}
type refreshReposDoneMsg struct {
	repos service.Repositories
}
type updateStatusMsg string
