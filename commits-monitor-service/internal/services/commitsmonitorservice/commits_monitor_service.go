package commitsmonitorservice

import (
	"commits-monitor-service/internal/constants"
	"commits-monitor-service/internal/constants/models"
	cmdsc "commits-monitor-service/internal/http/grpc/client/commits"
	rmdsc "commits-monitor-service/internal/http/grpc/client/repos"
	"commits-monitor-service/internal/http/grpc/protos/repos"
	"commits-monitor-service/internal/message-broker/rabbitmq"
	"commits-monitor-service/internal/pkg/githubrestclient"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CommentMonitorService struct {
	GithubRestClient             githubrestclient.GithubRestClient
	ReposMetaDataServiceClient   rmdsc.ReposMetaDataServiceClient
	CommitsMetaDataServiceClient cmdsc.CommitsMetaDataServiceClient
	Rabbit                       *amqp.Connection
}

func NewCommentMonitorService(
	githubRestClient githubrestclient.GithubRestClient,
	reposMetaDataServiceClient rmdsc.ReposMetaDataServiceClient,
	commitsMetaDataServiceClient cmdsc.CommitsMetaDataServiceClient,
	rabbit *amqp.Connection,
) CommentMonitorService {
	return CommentMonitorService{
		GithubRestClient:             githubRestClient,
		ReposMetaDataServiceClient:   reposMetaDataServiceClient,
		CommitsMetaDataServiceClient: commitsMetaDataServiceClient,
		Rabbit:                       rabbit,
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
	repositories, err := sc.ReposMetaDataServiceClient.GetRepositories()
	if err != nil {
		fmt.Println("CM Error getting repositories")
		fmt.Println("CM ERR:", err)
		return
	}

	var wg sync.WaitGroup
	for _, repo := range repositories {
		wg.Add(1)
		go func(repo *repos.Repository) {
			defer wg.Done()
			sc.fetchAndSaveCommitsForRepo(repo, fetchTime)
		}(repo)
	}
	wg.Wait()
}

func (sc *CommentMonitorService) fetchAndSaveCommitsForRepo(repo *repos.Repository, fetchTime time.Time) {
	lastFetchTime, _ := sc.CommitsMetaDataServiceClient.GetRepoLastFetchTime(repo.Name)
	since := ""
	if !(lastFetchTime == "") {
		since = lastFetchTime
	}

	commits, err := sc.GithubRestClient.FetchCommits(repo.Name, since)
	if err != nil {
		fmt.Println("CM Error fetching commits of ", repo.Name)
		fmt.Println("CM ERR:", err)
		return
	}

	fmt.Println(repo, " total commits ->", len(commits))
	sc.pushToQueue(repo.Name, fetchTime, commits)

}

// pushToQueue pushes a message into RabbitMQ
func (sc *CommentMonitorService) pushToQueue(repoName string, fetchTime time.Time, commits []models.CommitResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&event.Payload{
		Name: "commits",
		Data: struct {
			Repository string
			FetchTime  time.Time
			Commits    []models.CommitResponse
		}{
			Repository: repoName,
			FetchTime:  fetchTime,
			Commits:    commits,
		},
	}, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), constants.COMMITS_EVENT)
	if err != nil {
		return err
	}
	return nil
}
