package db_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"commits-manager-service/internal/constants/models"
	"github.com/google/uuid"

)



func createRandomRepository() models.Repository {
	return models.Repository{
		Name:            "test-repo-" + uuid.New().String(),
		Description:     "Test description",
		URL:             "https://github.com/test/test-repo",
		Language:        "Go",
		ForksCount:      10,
		StarsCount:      20,
		OpenIssuesCount: 1,
		WatchersCount:   5,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func TestInsertRepository(t *testing.T) {
	repo := createRandomRepository()

	insertedName, err := repositoryQueries.InsertRepository(repo)
	require.NoError(t, err)
	require.Equal(t, repo.Name, insertedName)

	retrievedRepo, err := repositoryQueries.GetRepositoryByName(repo.Name)
	require.NoError(t, err)
	require.Equal(t, repo.Name, retrievedRepo.Name)
	require.Equal(t, repo.Description, retrievedRepo.Description)

	repositoryQueries.DeleteRepository(retrievedRepo.Name)
}

func TestGetAllRepositories(t *testing.T) {
	repo1 := createRandomRepository()
	repo2 := createRandomRepository()

	_, err := repositoryQueries.InsertRepository(repo1)
	require.NoError(t, err)
	_, err = repositoryQueries.InsertRepository(repo2)
	require.NoError(t, err)

	repos, err := repositoryQueries.GetAllRepositories(1000000,0)
	require.NoError(t, err)
	require.Len(t, repos, 2)

	repositoryQueries.DeleteRepository(repo1.Name)
	repositoryQueries.DeleteRepository(repo2.Name)

}

func TestUpdateRepository(t *testing.T) {
	repo := createRandomRepository()

	_, err := repositoryQueries.InsertRepository(repo)
	require.NoError(t, err)

	repo.Description = "Updated description"
	err = repositoryQueries.UpdateRepository(repo)
	require.NoError(t, err)

	updatedRepo, err := repositoryQueries.GetRepositoryByName(repo.Name)
	require.NoError(t, err)
	require.Equal(t, "Updated description", updatedRepo.Description)
	repositoryQueries.DeleteRepository(repo.Name)
}

func TestDeleteRepository(t *testing.T) {
	repo := createRandomRepository()

	_, err := repositoryQueries.InsertRepository(repo)
	require.NoError(t, err)

	err = repositoryQueries.DeleteRepository(repo.Name)
	require.NoError(t, err)

	deletedRepo, err := repositoryQueries.GetRepositoryByName(repo.Name)
	require.Error(t, err)
	require.Nil(t, deletedRepo)
}

func TestSaveAllRepositories(t *testing.T) {
	repo1 := createRandomRepository()
	repo2 := createRandomRepository()

	err := repositoryQueries.SaveAllRepositories([]models.Repository{repo1, repo2})
	require.NoError(t, err)

	repos, err := repositoryQueries.GetAllRepositories(100000000,0)
	require.NoError(t, err)
	require.Len(t, repos, 2)

	// Update repositories
	repo1.Description = "Updated description 1"
	repo2.Description = "Updated description 2"

	err = repositoryQueries.SaveAllRepositories([]models.Repository{repo1, repo2})
	require.NoError(t, err)

	updatedRepo1, err := repositoryQueries.GetRepositoryByName(repo1.Name)
	require.NoError(t, err)
	require.Equal(t, "Updated description 1", updatedRepo1.Description)

	updatedRepo2, err := repositoryQueries.GetRepositoryByName(repo2.Name)
	require.NoError(t, err)
	require.Equal(t, "Updated description 2", updatedRepo2.Description)
}
