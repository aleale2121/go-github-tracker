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

type ReposDiscoveryService struct {
	GithubRestClient           githubrestclient.GithubRestClient
	ReposMetaDataServiceClient rmdsc.ReposMetaDataServiceClient
	Rabbit                     *amqp.Connection
}

func NewReposDiscoveryService(
	githubRestClient githubrestclient.GithubRestClient,
	reposMetaDataServiceClient rmdsc.ReposMetaDataServiceClient,
	rabbit *amqp.Connection,
) ReposDiscoveryService {
	return ReposDiscoveryService{
		GithubRestClient:           githubRestClient,
		ReposMetaDataServiceClient: reposMetaDataServiceClient,
		Rabbit:                     rabbit,
	}
}

func (sc *ReposDiscoveryService) ScheduleFetchingRepository(interval time.Duration) {
	log.Println("Fetching Repositories Started ")
	sc.fetchAndSaveRepositories() //Initial Fetch
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		go sc.fetchAndSaveRepositories()
	}
}

func (sc *ReposDiscoveryService) fetchAndSaveRepositories() {
	fetchTime := time.Now()
	since, err := sc.ReposMetaDataServiceClient.GetReposLastFetchTime()
	log.Println("last fetch time: ", since)
	if err != nil {
		log.Println("Error getting all repositories last fetch time")
		log.Println("ERR:", err)
	}

	if since == "0001-01-01T00:00:00Z" {
		since = ""
	}

	repositories, err := sc.GithubRestClient.FetchRepositories(since)
	if err != nil {
		log.Println("Error fetching repositories ")
		log.Println("ERR:", err)
		return
	}

	log.Println("total fetched repos: ", len(repositories))
	sc.pushToQueue(fetchTime, repositories)
}

// pushToQueue pushes a message into RabbitMQ
func (sc *ReposDiscoveryService) pushToQueue(fetchTime time.Time, repos []models.RepositoryResponse) error {
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

type ReposMetaData struct {
	FetchTime time.Time
	Repos     []models.RepositoryResponse
}
