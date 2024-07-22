package scheduler

import (
	"fmt"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/storage/db"
	"sync"
	"time"
)

const layout = "2006-01-02T15:04:05Z"

type SchedulerService struct {
	RepositoryPersistence db.RepositoryPersistence
	CommitPersistence     db.CommitPersistence
	MetadataPersistence   db.MetadataPersistence
	GithubRestClient      githubrestclient.GithubRestClient
}

func NewSchedulerService(repositoryPersistence db.RepositoryPersistence,
	commitPersistence db.CommitPersistence,
	MetadataPersistence db.MetadataPersistence,
	GithubRestClient githubrestclient.GithubRestClient,
) SchedulerService {
	return SchedulerService{
		RepositoryPersistence: repositoryPersistence,
		CommitPersistence:     commitPersistence,
		MetadataPersistence:   MetadataPersistence,
		GithubRestClient:      GithubRestClient,
	}
}

func (sc *SchedulerService) ScheduleFetchingRepository(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("Fetching Repositories Started ")
	for range ticker.C {
		go sc.fetchAndSaveRepositories()
	}
}

func (sc *SchedulerService) fetchAndSaveRepositories() {
	fetchTime := time.Now()
	lastFetchTime, _ := sc.MetadataPersistence.GetLastReposFetchTime()
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(layout)
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

func (sc *SchedulerService) ScheduleFetchingCommits(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("CM Fetching Commits Started ")
	for range ticker.C {
		go sc.fetchAndSaveCommits()
	}
}

func (sc *SchedulerService) fetchAndSaveCommits() {
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

func (sc *SchedulerService) fetchAndSaveCommitsForRepo(repo models.Repository, fetchTime time.Time) {
	lastFetchTime, _ := sc.MetadataPersistence.GetLastCommitFetchTime(repo.Name)
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(layout)
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
