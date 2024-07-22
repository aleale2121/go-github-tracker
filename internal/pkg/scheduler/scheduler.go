package scheduler

import (
	"fmt"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/storage/db"
	"time"
)

const layout = "2006-01-02T15:04:05-0700"

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
			continue
		}
		repositories := make([]models.Repository, 0)
		for _, repo := range githubRepositories {
			repositories = append(repositories, ConvertRepositoryResponseToRepository(repo))
		}

		err = sc.RepositoryPersistence.SaveAllRepositories(repositories)
		if err != nil {
			fmt.Println("Error saving repositories ")
			fmt.Println("ERR:", err)
			continue
		}
		if len(repositories) > 0 {
			err = sc.MetadataPersistence.SaveFetchReposMetadata(models.FetchReposMetadata{
				FetchedAt: fetchTime,
				Total:     len(repositories),
			})

			if err != nil {
				fmt.Println("Error updating last fetch time ")
				fmt.Println("ERR:", err)
				continue
			}
		}

	}
}

func (sc *SchedulerService) ScheduleFetchingCommits(interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("CM Fetching Commits Started ")
	for range ticker.C {
		fetchTime := time.Now()
		repositories, err := sc.RepositoryPersistence.GetAllRepositories()
		if err != nil {
			fmt.Println("CM Error getting repositories")
			fmt.Println("CM ERR:", err)
			continue
		}
		for _, repo := range repositories {
			lastFetchTime, _ := sc.MetadataPersistence.GetLastCommitFetchTime(repo.Name)
			since := ""
			if !lastFetchTime.IsZero() {
				since = lastFetchTime.UTC().Format(layout)
			}

			fetchedCommits, err := sc.GithubRestClient.FetchCommits(repo.Name, since)
			if err != nil {
				fmt.Println("CM Error fetching commits of ", repo.Name)
				fmt.Println("CM ERR:", err)
				continue
			}
			commits := make([]models.Commit, 0)
			for _, commit := range fetchedCommits {
				commits = append(commits, ConvertCommitResponseToCommit(commit, repo.Name))
			}
			err = sc.CommitPersistence.SaveAllCommits(commits)
			if err != nil {
				fmt.Println("CM Error saving commits of ", repo.Name)
				fmt.Println("CM ERR:", err)
				continue
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
					continue
				}
			}
		}

	}
}
