package db_test

import (
	"commits-manager-service/internal/constants/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomCommit(t *testing.T, repoName string) models.Commit {
	commit := models.Commit{
		SHA:            uuid.New().String(),
		URL:            "http://example.com/commit",
		Message:        "Test commit message",
		AuthorName:     "Author",
		AuthorDate:     time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		RepositoryName: repoName,
	}

	err := commitsQueries.InsertCommit(commit)
	require.NoError(t, err)
	return commit
}

func TestInsertCommit(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit:=createRandomCommit(t, repoName)
	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestGetCommitBySHA(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	retrievedCommit, err := commitsQueries.GetCommitBySHA(commit.SHA)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedCommit)
	require.Equal(t, commit.SHA, retrievedCommit.SHA)
	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestUpdateCommit(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	commit.Message = "Updated commit message"
	err = commitsQueries.UpdateCommit(commit)
	require.NoError(t, err)

	retrievedCommit, err := commitsQueries.GetCommitBySHA(commit.SHA)
	require.NoError(t, err)
	require.NotEmpty(t, retrievedCommit)
	require.Equal(t, "Updated commit message", retrievedCommit.Message)

	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestDeleteCommit(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	err = commitsQueries.DeleteCommit(commit.SHA)
	require.NoError(t, err)

	retrievedCommit, err := commitsQueries.GetCommitBySHA(commit.SHA)
	require.Error(t, err)
	require.Empty(t, retrievedCommit)

	repositoryQueries.DeleteRepository(repoName)

}

func TestGetAllCommits(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit1 := createRandomCommit(t, repoName)
	commit2 := createRandomCommit(t, repoName)

	commits, err := commitsQueries.GetAllCommits()
	require.NoError(t, err)
	require.NotEmpty(t, commits)
	require.Len(t, commits, 2)
	require.Equal(t, commit1.SHA, commits[0].SHA)
	require.Equal(t, commit2.SHA, commits[1].SHA)
	commitsQueries.DeleteCommit(commit1.SHA)
	commitsQueries.DeleteCommit(commit2.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestGetCommitsByRepoName(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	commits, err := commitsQueries.GetCommitsByRepoName(repoName)
	require.NoError(t, err)
	require.NotEmpty(t, commits)
	require.Equal(t, commit.RepositoryName, commits[0].RepositoryName)
	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestGetTopCommitAuthors(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	authors, err := commitsQueries.GetTopCommitAuthors(1)
	require.NoError(t, err)
	require.NotEmpty(t, authors)
	require.Equal(t, commit.AuthorName, authors[0].Name)
	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)
}

func TestGetTopCommitAuthorsByRepo(t *testing.T) {
	repoName := uuid.New().String()
	_, err := repositoryQueries.InsertRepository(models.Repository{
		Name:            repoName,
		Description:     "Test Repository",
		URL:             "http://example.com/repo",
		Language:        "Go",
		ForksCount:      1,
		StarsCount:      1,
		OpenIssuesCount: 0,
		WatchersCount:   1,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	commit := createRandomCommit(t, repoName)
	authors, err := commitsQueries.GetTopCommitAuthorsByRepo(repoName, 1)
	require.NoError(t, err)
	require.NotEmpty(t, authors)
	require.Equal(t, commit.AuthorName, authors[0].Name)
	commitsQueries.DeleteCommit(commit.SHA)
	repositoryQueries.DeleteRepository(repoName)

}
