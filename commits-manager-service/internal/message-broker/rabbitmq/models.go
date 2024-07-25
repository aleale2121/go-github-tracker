package event

import (
	"commits-manager-service/internal/constants/models"
	"time"
)

type CommitMetaData struct {
	Repository string
	FetchTime  time.Time
	Commits    []models.CommitResponse
}

type ReposMetaData struct {
	FetchTime time.Time
	Repos     []models.RepositoryResponse
}
