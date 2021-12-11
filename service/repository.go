package service

type Repository struct{}

type Repositories []Repository

type RepositoryService struct {
}

func (r *RepositoryService) GetRepositoriesByOwners(owners []string) Repositories {
	res := Repositories{}
	return res
}
