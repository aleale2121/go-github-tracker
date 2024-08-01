package repos

import "commits-manager-service/internal/storage/db"
import "commits-manager-service/internal/constants/models"

type RepositoryManagerService struct {
	RepositoryPersistence db.GitReposRepository
}

func NewRepositoryManagerService(repositoryPersistence db.GitReposRepository) RepositoryManagerService {
	return RepositoryManagerService{RepositoryPersistence: repositoryPersistence}
}

func (rc RepositoryManagerService) GetRepositories(limit, offset int) ([]*models.Repository, error) {
	return rc.RepositoryPersistence.GetAllRepositories(limit, offset)
}
