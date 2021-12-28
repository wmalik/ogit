package browser

import "ogit/service"

type refreshReposMsg struct{}
type refreshReposDoneMsg struct {
	repos      service.Repositories
	rateLimits string
}
type updateStatusMsg string
