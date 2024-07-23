package commentmanagerservice

import "go-github-tracker/internal/storage/db"
import "go-github-tracker/internal/constants/models"

type CommitsManagerService struct {
	CommitsPersistence db.CommitPersistence
}

func NewCommitsManagerService(commitsPersistence db.CommitPersistence) CommitsManagerService {
	return CommitsManagerService{CommitsPersistence: commitsPersistence}
}

func (rc CommitsManagerService) GetCommitsByRepositoryName(repoName string) ([]*models.Commit, error) {
	return rc.CommitsPersistence.GetCommitsByRepoName(repoName)
}

func (rc CommitsManagerService) GetTopCommitAuthors(limit int) ([]*models.CommitAuthor, error) {
	return rc.CommitsPersistence.GetTopCommitAuthors(limit)
}
func (rc CommitsManagerService) GetTopCommitAuthorsByRepoName(repoName string, limit int) ([]*models.CommitAuthor, error) {
	return rc.CommitsPersistence.GetTopCommitAuthorsByRepo(repoName, limit)
}
