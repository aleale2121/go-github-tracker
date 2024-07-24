package reposdiscoveryservice

import (
	"encoding/json"
	"fmt"
	"repos-discovery-service/internal/constants"
	"repos-discovery-service/internal/constants/models"
	"repos-discovery-service/internal/message-broker/rabbitmq"
	"repos-discovery-service/internal/pkg/githubrestclient"

	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ReposDiscoveryService struct {
	GithubRestClient githubrestclient.GithubRestClient
	Rabbit           *amqp.Connection
}

func NewReposDiscoveryService(
	githubRestClient githubrestclient.GithubRestClient,
	rabbit *amqp.Connection,
) ReposDiscoveryService {
	return ReposDiscoveryService{
		GithubRestClient: githubRestClient,
		Rabbit:           rabbit,
	}
}

func (sc *ReposDiscoveryService) ScheduleFetchingRepository(interval time.Duration) {
	sc.fetchAndSaveRepositories() //Initial Fetch
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	fmt.Println("Fetching Repositories Started ")
	for range ticker.C {
		go sc.fetchAndSaveRepositories()
	}
}

func (sc *ReposDiscoveryService) fetchAndSaveRepositories() {
	lastFetchTime := time.Date(2024, 7, 1, 1, 1, 1, 1, time.Local)
	since := ""
	if !lastFetchTime.IsZero() {
		since = lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT)
	}

	repositories, err := sc.GithubRestClient.FetchRepositories(since)
	if err != nil {
		fmt.Println("Error fetching repositories ")
		fmt.Println("ERR:", err)
		return
	}

	fmt.Println(repositories)
	sc.pushToQueue(repositories)

	// repositories := make([]models.Repository, len(githubRepositories))
	// for i, repo := range githubRepositories {
	// 	repositories[i] = ConvertRepositoryResponseToRepository(repo)
	// }

	// err = sc.RepositoryPersistence.SaveAllRepositories(repositories)
	// if err != nil {
	// 	fmt.Println("Error saving repositories ")
	// 	fmt.Println("ERR:", err)
	// 	return
	// }
	// if len(repositories) > 0 {
	// 	err = sc.MetadataPersistence.SaveFetchReposMetadata(models.FetchReposMetadata{
	// 		FetchedAt: fetchTime,
	// 		Total:     len(repositories),
	// 	})

	// 	if err != nil {
	// 		fmt.Println("Error updating last fetch time ")
	// 		fmt.Println("ERR:", err)
	// 		return
	// 	}
	// }
}

// pushToQueue pushes a message into RabbitMQ
func (sc *ReposDiscoveryService) pushToQueue(repos []models.RepositoryResponse) error {
	emitter, err := event.NewEventEmitter(sc.Rabbit)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(&repos, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), constants.REPOS_FETCHED)
	if err != nil {
		return err
	}
	return nil
}

// func ConvertRepositoryResponseToRepository(response models.RepositoryResponse) models.Repository {
// 	description := ""
// 	if response.Description != nil {
// 		description = response.Description.(string)
// 	}

// 	return models.Repository{
// 		Name:            response.Name,
// 		Description:     description,
// 		URL:             response.HTMLURL,
// 		Language:        response.Language,
// 		ForksCount:      response.ForksCount,
// 		StarsCount:      response.StargazersCount,
// 		OpenIssuesCount: response.OpenIssuesCount,
// 		WatchersCount:   response.WatchersCount,
// 		CreatedAt:       response.CreatedAt,
// 		UpdatedAt:       response.UpdatedAt,
// 	}
// }
