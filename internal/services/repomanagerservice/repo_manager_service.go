package repomanagerservice

import "go-github-tracker/internal/storage/db"
import "go-github-tracker/internal/constants/models"

type RepositoryManagerService struct {
	RepositoryPersistence db.RepositoryPersistence
}

func NewRepositoryManagerService(repositoryPersistence db.RepositoryPersistence) RepositoryManagerService {
	return RepositoryManagerService{RepositoryPersistence: repositoryPersistence}
}

func (rc RepositoryManagerService) GetRepositories() ([]*models.Repository, error) {
	return rc.RepositoryPersistence.GetAllRepositories()
}
