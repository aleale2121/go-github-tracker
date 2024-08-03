package commitsmonitorservice

import (
	"commits-monitor-service/internal/constants"
	"commits-monitor-service/internal/constants/models"
	cmdsc "commits-monitor-service/internal/http/grpc/client/commits"
	rmdsc "commits-monitor-service/internal/http/grpc/client/repos"
	"commits-monitor-service/internal/message-broker/rabbitmq"
	"commits-monitor-service/internal/pkg/githubrestclient"
	"encoding/json"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const perPage = 100

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
	sc.fetchAndSaveCommits()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		go sc.fetchAndSaveCommits()
	}
}

func (sc *CommentMonitorService) fetchAndSaveCommits() {
	repositories, err := sc.ReposMetaDataServiceClient.GetRepositoryNames()
	if err != nil {
		log.Println("CMOS: error getting repository Names")
		log.Println("CMOS: err:", err)
		return
	}
	log.Printf("CMOS: fetching commits of %d repositories started\n", len(repositories))
	var wg sync.WaitGroup
	for _, repo := range repositories {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			sc.fetchAndSaveCommitsForRepo(repo)
		}(repo)
	}
	wg.Wait()
}

func (sc *CommentMonitorService) fetchAndSaveCommitsForRepo(repo string) {
	commitFetchHistory, err := sc.CommitsMetaDataServiceClient.GetCommitFetchHistory(repo)
	if err != nil {
		log.Println("CMOS: error getting a repository last commit fetch history")
		log.Println("CMOS: err:", err)
	}

	page := commitFetchHistory.LastPage + 1

	log.Printf("CMOS: fetching commits of <%s> started from page=> %d\n", repo, page)

	var totalCommitsFetched int

	for {
		commits, err := sc.GithubRestClient.FetchCommits(repo, perPage, page)
		if err != nil {
			log.Println("CMOS: error fetching commits of ", repo)
			log.Println("CMOS: err:", err)
			return
		}

		if len(commits) == 0 {
			break
		}
        
		log.Printf("CMOS: pulled %d commits %s \n", len(commits), repo)
		
		fetchTime := commits[len(commits)-1].Commit.Author.Date
		sc.pushToQueue(repo, fetchTime, page, commits)

		totalCommitsFetched += len(commits)
		page++
	}

	log.Printf("CMOS: repo <%s>  total commits: %d pulled\n", repo, totalCommitsFetched)
}

// pushToQueue pushes a message into RabbitMQ
func (sc *CommentMonitorService) pushToQueue(repoName string, fetchTime time.Time, lastPage int32, commits []models.CommitResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&event.Payload{
		Name: "commits",
		Data: CommitMetaData{
			Repository: repoName,
			FetchTime:  fetchTime,
			Commits:    commits,
			LastPage:   int(lastPage),
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

type CommitMetaData struct {
	LastPage   int
	Repository string
	FetchTime  time.Time
	Commits    []models.CommitResponse
}
