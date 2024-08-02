package reposdiscoveryservice

import (
	"encoding/json"
	"log"
	"repos-discovery-service/internal/constants"
	"repos-discovery-service/internal/constants/models"
	rmdsc "repos-discovery-service/internal/http/grpc/client/repos"
	"repos-discovery-service/internal/message-broker/rabbitmq"
	"repos-discovery-service/internal/pkg/githubrestclient"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const perPage = 30

type ReposDiscoveryService struct {
	GithubRestClient           githubrestclient.GithubRestClient
	ReposMetaDataServiceClient rmdsc.RepositoriesServiceClient
	Rabbit                     *amqp.Connection
}

func NewReposDiscoveryService(
	githubRestClient githubrestclient.GithubRestClient,
	reposMetaDataServiceClient rmdsc.RepositoriesServiceClient,
	rabbit *amqp.Connection,
) ReposDiscoveryService {
	return ReposDiscoveryService{
		GithubRestClient:           githubRestClient,
		ReposMetaDataServiceClient: reposMetaDataServiceClient,
		Rabbit:                     rabbit,
	}
}

func (sc *ReposDiscoveryService) ScheduleDiscoveringNewRepository(interval time.Duration) {
	log.Println("RDS: discovering New Repositories Started ")
	sc.discoverAndSaveNewRepositories()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		go sc.discoverAndSaveNewRepositories()
	}
}

func (sc *ReposDiscoveryService) ScheduleFetchingRepositoryMetadata(interval time.Duration) {
	log.Println("RDS: fetching Repositories Metadata Started ")
	sc.fetchRepositoriesMetadata()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		go sc.fetchRepositoriesMetadata()
	}
}

func (sc *ReposDiscoveryService) discoverAndSaveNewRepositories() {
	fetchTime := time.Now()
	since, err := sc.ReposMetaDataServiceClient.GetReposLastFetchTime()
	log.Println("RDS: last discovering repositories fetch time: ", since)
	if err != nil {
		log.Println("RDS: Error getting all repositories last fetch time")
		log.Println("RDS: ERR:", err)
	}

	if since == "0001-01-01T00:00:00Z" {
		since = ""
	}

	page := 1
	var totalRepositories int

	for {
		repositories, err := sc.GithubRestClient.FetchRepositories(since, perPage, page)
		if err != nil {
			log.Println("RDS: error fetching repositories ")
			log.Println("RDS: err:", err)
			return
		}

		if len(repositories) == 0 {
			break
		}
		totalRepositories += len(repositories)
		page++
		log.Printf("RDS: pulled %d repositories \n", len(repositories))
		sc.pushNewRepositoriesToQueue(fetchTime, repositories)

	}

	log.Println("RDS: total fetched repos: ", totalRepositories)
}

func (sc *ReposDiscoveryService) fetchRepositoriesMetadata() {
	repositories, err := sc.ReposMetaDataServiceClient.GetRepositoryNames()
	if err != nil {
		log.Println("RDS: error getting repository names")
		log.Println("RDS: err: ", err)
	}

	for _, repoName := range repositories {
		repository, err := sc.GithubRestClient.FetchRepositoryMetadata(repoName)
		if err != nil {
			log.Println("RDS: error getting repository meta data")
			log.Println("RDS: err:", err)
		}
		sc.pushRepositoryMetaDataToQueue(repository)
	}
}

// pushNewRepositoriesToQueue pushes a message into RabbitMQ
func (sc *ReposDiscoveryService) pushNewRepositoriesToQueue(fetchTime time.Time, repos []models.RepositoryResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&event.Payload{
		Name: "repos",
		Data: ReposMetaData{
			FetchTime: fetchTime,
			Repos:     repos,
		},
	}, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), constants.REPOS_EVENT)
	if err != nil {
		return err
	}
	return nil
}

// pushRepositoryMetaDataToQueue pushes a message into RabbitMQ
func (sc *ReposDiscoveryService) pushRepositoryMetaDataToQueue(repo models.RepositoryResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&event.Payload{
		Name: "repo",
		Data: repo,
	}, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), constants.REPOS_EVENT)
	if err != nil {
		return err
	}
	return nil
}

type ReposMetaData struct {
	FetchTime time.Time
	Repos     []models.RepositoryResponse
}
