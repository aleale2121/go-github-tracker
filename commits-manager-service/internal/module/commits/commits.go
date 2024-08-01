package commits

import (
	"commits-manager-service/internal/constants/models"
	"commits-manager-service/internal/storage/db"
	"time"
)

type CommitsManagerService struct {
	CommitsPersistence db.CommitRepository
}

func NewCommitsManagerService(commitsPersistence db.CommitRepository) CommitsManagerService {
	return CommitsManagerService{CommitsPersistence: commitsPersistence}
}

func (rc CommitsManagerService) GetCommitsByRepositoryName(repoName string, limit, offset int, startDate, endDate time.Time) ([]*models.Commit, error) {
	return rc.CommitsPersistence.GetCommitsByRepoName(repoName, limit, offset, startDate, endDate)
}

func (rc CommitsManagerService) GetTopCommitAuthors(limit int) ([]*models.CommitAuthor, error) {
	return rc.CommitsPersistence.GetTopCommitAuthors(limit)
}
func (rc CommitsManagerService) GetTopCommitAuthorsByRepoName(repoName string, limit int) ([]*models.CommitAuthor, error) {
	return rc.CommitsPersistence.GetTopCommitAuthorsByRepo(repoName, limit)
}

func (rc CommitsManagerService) GetTotalCommitsByRepositoryName(repoName string, startDate, endDate time.Time) (int, error) {
	return rc.CommitsPersistence.GetTotalCommitsByRepoName(repoName, startDate, endDate)
}
