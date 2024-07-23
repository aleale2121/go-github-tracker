package repodiscoveryservice

import (
	"fmt"
	"go-github-tracker/internal/constants"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/storage/db"
	"time"
)

type CommentDiscoveryService struct {
	RepositoryPersistence db.RepositoryPersistence
	CommitPersistence     db.CommitPersistence
	MetadataPersistence   db.MetadataPersistence
	GithubRestClient      githubrestclient.GithubRestClient
}

func NewCommentDiscoveryService(repositoryPersistence db.RepositoryPersistence,
	commitPersistence db.CommitPersistence,
	MetadataPersistence db.MetadataPersistence,
	GithubRestClient githubrestclient.GithubRestClient,
) CommentDiscoveryService {
	return CommentDiscoveryService{
		RepositoryPersistence: repositoryPersistence,
		CommitPersistence:     commitPersistence,
		MetadataPersistence:   MetadataPersistence,
		GithubRestClient:      GithubRestClient,
	}
}

func (sc *CommentDiscoveryService) ScheduleFetchingRepository(interval time.Duration) {
	sc.fetchAndSaveRepositories() //Initial Fetch
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("Fetching Repositories Started ")
	for range ticker.C {
		go sc.fetchAndSaveRepositories()
	}
}

func (sc *CommentDiscoveryService) fetchAndSaveRepositories() {
	fetchTime := time.Now()
	lastFetchTime, _ := sc.MetadataPersistence.GetLastReposFetchTime()
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT)
	}

	githubRepositories, err := sc.GithubRestClient.FetchRepositories(since)
	if err != nil {
		fmt.Println("Error fetching repositories ")
		fmt.Println("ERR:", err)
		return
	}
	
	repositories := make([]models.Repository, len(githubRepositories))
	for i, repo := range githubRepositories {
		repositories[i] = ConvertRepositoryResponseToRepository(repo)
	}

	err = sc.RepositoryPersistence.SaveAllRepositories(repositories)
	if err != nil {
		fmt.Println("Error saving repositories ")
		fmt.Println("ERR:", err)
		return
	}
	if len(repositories) > 0 {
		err = sc.MetadataPersistence.SaveFetchReposMetadata(models.FetchReposMetadata{
			FetchedAt: fetchTime,
			Total:     len(repositories),
		})

		if err != nil {
			fmt.Println("Error updating last fetch time ")
			fmt.Println("ERR:", err)
			return
		}
	}
}

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
