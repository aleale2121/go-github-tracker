package commitsmonitorservice

import (
	"fmt"
	"go-github-tracker/internal/constants"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/storage/db"
	"sync"
	"time"
)


type CommentMonitorService struct {
	RepositoryPersistence db.RepositoryPersistence
	CommitPersistence     db.CommitPersistence
	MetadataPersistence   db.MetadataPersistence
	GithubRestClient      githubrestclient.GithubRestClient
}

func NewCommentMonitorService(repositoryPersistence db.RepositoryPersistence,
	commitPersistence db.CommitPersistence,
	MetadataPersistence db.MetadataPersistence,
	GithubRestClient githubrestclient.GithubRestClient,
) CommentMonitorService {
	return CommentMonitorService{
		RepositoryPersistence: repositoryPersistence,
		CommitPersistence:     commitPersistence,
		MetadataPersistence:   MetadataPersistence,
		GithubRestClient:      GithubRestClient,
	}
}

func (sc *CommentMonitorService) ScheduleFetchingCommits(interval time.Duration) {
	fmt.Println("Fetching Commits Started ")
	sc.fetchAndSaveCommits() 
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		go sc.fetchAndSaveCommits()
	}
}

func (sc *CommentMonitorService) fetchAndSaveCommits() {
	fetchTime := time.Now()
	repositories, err := sc.RepositoryPersistence.GetAllRepositories()
	if err != nil {
		fmt.Println("CM Error getting repositories")
		fmt.Println("CM ERR:", err)
		return
	}

	var wg sync.WaitGroup
	for _, repo := range repositories {
		wg.Add(1)
		go func(repo models.Repository) {
			defer wg.Done()
			sc.fetchAndSaveCommitsForRepo(repo, fetchTime)
		}(*repo)
	}
	wg.Wait()
}

func (sc *CommentMonitorService) fetchAndSaveCommitsForRepo(repo models.Repository, fetchTime time.Time) {
	lastFetchTime, _ := sc.MetadataPersistence.GetLastCommitFetchTime(repo.Name)
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT)
	}

	fetchedCommits, err := sc.GithubRestClient.FetchCommits(repo.Name, since)
	if err != nil {
		fmt.Println("CM Error fetching commits of ", repo.Name)
		fmt.Println("CM ERR:", err)
		return
	}
	commits := make([]models.Commit, len(fetchedCommits))
	for i, commit := range fetchedCommits {
		commits[i] = ConvertCommitResponseToCommit(commit, repo.Name)
	}

	err = sc.CommitPersistence.SaveAllCommits(commits)
	if err != nil {
		fmt.Println("CM Error saving commits of ", repo.Name)
		fmt.Println("CM ERR:", err)
		return
	}

	if len(commits) > 0 {
		err = sc.MetadataPersistence.SaveFetchCommitsMetadata(models.FetchCommitsMetadata{
			RepositoryName: repo.Name,
			FetchedAt:      fetchTime,
			Total:          len(commits),
		})

		if err != nil {
			fmt.Println("CM Error updating last commit fetch time ", repo.Name)
			fmt.Println("CM ERR:", err)
			return
		}
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
