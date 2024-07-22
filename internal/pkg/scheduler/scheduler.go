package scheduler

import (
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/storage/db"
	"time"
)

type SchedulerService struct {
	RepositoryPersistence db.RepositoryPersistence
	CommitPersistence     db.CommitPersistence
	GithubRestClient      githubrestclient.GithubRestClient
}

func NewSchedulerService(repositoryPersistence db.RepositoryPersistence,
	commitPersistence db.CommitPersistence) SchedulerService {
	return SchedulerService{
		RepositoryPersistence: repositoryPersistence,
		CommitPersistence:     commitPersistence,
	}
}

func (sc *SchedulerService) ScheduleFetchingRepository(wait chan bool, interval time.Duration) {
	ticker := time.NewTicker(interval)
    defer ticker.Stop()

	for {
        select {
        case <-ticker.C:
            githubRepositories, err := sc.GithubRestClient.FetchRepositories();
            if err != nil {
                continue
            }

			repositories:=make([]models.Repository,0)
			for _,repo:=range githubRepositories{
				repositories = append(repositories, ConvertRepositoryResponseToRepository(repo))
			}
        
     
            // Save repository data
            // saveRepository(db, repository)

            // Fetch and save commits
            // since := time.Now().Add(-interval).Format(time.RFC3339)
            // commits, err := fetchCommits(client, owner, repo, since)
            // if err != nil {
            //     continue
            // }

            // var repoRecord Repository
            // db.Where("name = ?", repo).First(&repoRecord)
            // saveCommits(db, repoRecord.ID, commits)

        }
    }
}
