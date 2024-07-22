package scheduler

import (
	"go-github-tracker/internal/constants/models"
)

func ConvertRepositoryResponseToRepository(response models.RepositoryResponse) models.Repository {
	description := ""
	if response.Description != nil {
		description = response.Description.(string)
	}

	return models.Repository{
		Name:            response.Name,
		Description:     description,
		URL:             response.HTMLURL,
		Language:        response.Language,
		ForksCount:      response.ForksCount,
		StarsCount:      response.StargazersCount,
		OpenIssuesCount: response.OpenIssuesCount,
		WatchersCount:   response.WatchersCount,
		CreatedAt:       response.CreatedAt,
		UpdatedAt:       response.UpdatedAt,
	}
}

func ConvertCommitResponseToCommit(response models.CommitResponse, repositoryName string) models.Commit {
	return models.Commit{
		SHA:            response.Sha,
		URL:            response.URL,
		Message:        response.Commit.Message,
		AuthorName:     response.Commit.Author.Name,
		AuthorDate:     response.Commit.Author.Date,
		CreatedAt:      response.Commit.Author.Date,
		UpdatedAt:      response.Commit.Committer.Date,
		RepositoryName: repositoryName,
	}
}
