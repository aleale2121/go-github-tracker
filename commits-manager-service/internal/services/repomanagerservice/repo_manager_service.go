package repomanagerservice

import "commits-manager-service/internal/storage/db"
import "commits-manager-service/internal/constants/models"

type RepositoryManagerService struct {
	RepositoryPersistence db.RepositoryPersistence
}

func NewRepositoryManagerService(repositoryPersistence db.RepositoryPersistence) RepositoryManagerService {
	return RepositoryManagerService{RepositoryPersistence: repositoryPersistence}
}

func (rc RepositoryManagerService) GetRepositories() ([]*models.Repository, error) {
	return rc.RepositoryPersistence.GetAllRepositories()
}
