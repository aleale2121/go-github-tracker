package event

import (
	"commits-manager-service/internal/constants/models"
	"time"
)

type CommitMetaData struct {
	LastPage int
	Repository string
	FetchTime  time.Time
	Commits    []models.CommitResponse
}

type ReposMetaData struct {
	LastPage  int
	FetchTime time.Time
	Repos     []models.RepositoryResponse
}
