package commitsmonitorservice

import (
	"commits-monitor-service/internal/constants"
	"commits-monitor-service/internal/constants/models"
	"commits-monitor-service/internal/message-broker/rabbitmq"
	"commits-monitor-service/internal/pkg/githubrestclient"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CommentMonitorService struct {
	GithubRestClient githubrestclient.GithubRestClient
	Rabbit           *amqp.Connection
}

func NewCommentMonitorService(
	githubRestClient githubrestclient.GithubRestClient,
	rabbit *amqp.Connection,
) CommentMonitorService {
	return CommentMonitorService{
		GithubRestClient: githubRestClient,
		Rabbit:           rabbit,
	}
}

func (sc *CommentMonitorService) ScheduleFetchingCommits(interval time.Duration) {
	// Wait 1 minute for first repository fetch
	timer := time.After(60 * time.Second)
	<-timer

	// Start Fetching Commits
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
	repositories := []*models.Repository{}
	// if err != nil {
	// 	fmt.Println("CM Error getting repositories")
	// 	fmt.Println("CM ERR:", err)
	// 	return
	// }

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
	// lastFetchTime, _ := sc.MetadataPersistence.GetLastCommitFetchTime(repo.Name)
	lastFetchTime := time.Date(2024, 7, 1, 1, 1, 1, 1, time.Local)
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT)
	}

	commits, err := sc.GithubRestClient.FetchCommits(repo.Name, since)
	if err != nil {
		fmt.Println("CM Error fetching commits of ", repo.Name)
		fmt.Println("CM ERR:", err)
		return
	}
	// commits := make([]models.Commit, len(fetchedCommits))
	// for i, commit := range fetchedCommits {
	// 	commits[i] = ConvertCommitResponseToCommit(commit, repo.Name)
	// }
	fmt.Println(commits)
	sc.pushToQueue(commits)
	// err = sc.CommitPersistence.SaveAllCommits(commits)
	// if err != nil {
	// 	fmt.Println("CM Error saving commits of ", repo.Name)
	// 	fmt.Println("CM ERR:", err)
	// 	return
	// }

	// if len(commits) > 0 {
	// 	err = sc.MetadataPersistence.SaveFetchCommitsMetadata(models.FetchCommitsMetadata{
	// 		RepositoryName: repo.Name,
	// 		FetchedAt:      fetchTime,
	// 		Total:          len(commits),
	// 	})

	// 	if err != nil {
	// 		fmt.Println("CM Error updating last commit fetch time ", repo.Name)
	// 		fmt.Println("CM ERR:", err)
	// 		return
	// 	}
	// }
}

// pushToQueue pushes a message into RabbitMQ
func (sc *CommentMonitorService) pushToQueue(commits []models.CommitResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&commits, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), constants.COMMITS_FETCHED)
	if err != nil {
		return err
	}
	return nil
}

// func ConvertCommitResponseToCommit(response models.CommitResponse, repositoryName string) models.Commit {
// 	return models.Commit{
// 		SHA:            response.Sha,
// 		URL:            response.URL,
// 		Message:        response.Commit.Message,
// 		AuthorName:     response.Commit.Author.Name,
// 		AuthorDate:     response.Commit.Author.Date,
// 		CreatedAt:      response.Commit.Author.Date,
// 		UpdatedAt:      response.Commit.Committer.Date,
// 		RepositoryName: repositoryName,
// 	}
// }
