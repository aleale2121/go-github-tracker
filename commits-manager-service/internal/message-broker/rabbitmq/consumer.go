package event

import (
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/constants/models"
	"commits-manager-service/internal/storage/db"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn                  *amqp.Connection
	queueName             string
	CommitPersistence     db.CommitRepository
	RepositoryPersistence db.GitReposRepository
}

func NewConsumer(conn *amqp.Connection, queueName string,
	commitPersistence db.CommitRepository,
	repositoryPersistence db.GitReposRepository) (Consumer, error) {
	consumer := Consumer{
		conn:                  conn,
		queueName:             queueName,
		CommitPersistence:     commitPersistence,
		RepositoryPersistence: repositoryPersistence,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			constants.GITHUB_API_TOPIC,
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			switch payload.Name {
			case "repos":
				go consumer.proccessAndSaveNewRepos(payload)
			case "repo":
				go consumer.proccessAndUpdateRepoMetaData(payload)
			case "commits":
				go consumer.proccessAndSaveCommits(payload)
			default:
				log.Println("recieved payload-->", payload)
			}
		}
	}()

	fmt.Printf("Consumer: Waiting for message [Exchange, Queue] [github_api_topics, %s]\n", q.Name)
	<-forever

	return nil
}

func (consumer *Consumer) proccessAndSaveCommits(entry Payload) {
	jsonData, _ := json.MarshalIndent(entry.Data, "", "\t")

	var commitMetaData CommitMetaData
	err := json.Unmarshal(jsonData, &commitMetaData)

	if err == nil {
		log.Println("Consumer-Recieved-Commit->", commitMetaData.Repository, len(commitMetaData.Commits))
		if len(commitMetaData.Commits) > 0 {
			commits := make([]models.Commit, len(commitMetaData.Commits))
			for i, commit := range commitMetaData.Commits {
				commits[i] = ConvertCommitResponseToCommit(commit, commitMetaData.Repository)
			}

			err := consumer.CommitPersistence.SaveAllCommits(commits)
			if err != nil {
				fmt.Println("Consumer: Error saving commits of ", commitMetaData.Repository)
				fmt.Println("Consumer: ERR:", err)
				return
			}

			err = consumer.CommitPersistence.SaveCommitsFetchData(models.CommitsFetchData{
				RepositoryName: commitMetaData.Repository,
				FetchedAt:      commitMetaData.FetchTime,
				Total:          len(commits),
			})

			if err != nil {
				fmt.Println("Consumer: Error updating last commit fetch time ", commitMetaData.Repository)
				fmt.Println("Consumer: ERR:", err)
				return
			}
		}
	} else {
		log.Println("Consumer: Cannot Convert To Commit MetaData")
	}

}

func (consumer *Consumer) proccessAndSaveNewRepos(entry Payload) {
	jsonData, _ := json.MarshalIndent(entry.Data, "", "\t")

	var reposMetaData ReposMetaData
	err := json.Unmarshal(jsonData, &reposMetaData)

	if err == nil {
		log.Println("Consumer-Recieved-Repositories->", len(reposMetaData.Repos))
		if len(reposMetaData.Repos) > 0 {
			repositories := make([]models.Repository, len(reposMetaData.Repos))
			for i, repo := range reposMetaData.Repos {
				repositories[i] = ConvertRepositoryResponseToRepository(repo)
			}

			err := consumer.RepositoryPersistence.SaveAllRepositories(repositories)
			if err != nil {
				fmt.Println("Consumer: Error saving repositories ")
				fmt.Println("Consumer: ERR:", err)
				return
			}
			err = consumer.RepositoryPersistence.SaveReposFetchData(models.ReposFetchData{
				FetchedAt: reposMetaData.FetchTime,
				Total:     len(repositories),
			})

			if err != nil {
				fmt.Println("Consumer: Error updating last fetch time ")
				fmt.Println("Consumer: ERR:", err)
				return
			}

		}
	} else {
		log.Println("Consumer: Cannot Convert To RepositoryMetaData")
	}
}

func (consumer *Consumer) proccessAndUpdateRepoMetaData(entry Payload) {
	jsonData, _ := json.MarshalIndent(entry.Data, "", "\t")

	var repository models.RepositoryResponse
	err := json.Unmarshal(jsonData, &repository)

	if err == nil {
		log.Println("Consumer-Recieved-Repository MetaData->", repository.Name)

		err := consumer.RepositoryPersistence.UpdateRepository(
			ConvertRepositoryResponseToRepository(repository))
		if err != nil {
			fmt.Println("Consumer: Error updating repository metadat")
			fmt.Println("Consumer: ERR:", err)
			return
		}

	} else {
		log.Println("Consumer: Cannot Convert To Repository")
	}
}

func ConvertCommitResponseToCommit(response models.CommitResponse, repositoryName string) models.Commit {
	return models.Commit{
		SHA:            response.Sha,
		URL:            response.URL,
		Message:        response.Commit.Message,
		AuthorName:     response.Commit.Author.Name,
		AuthorDate:     response.Commit.Author.Date,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		RepositoryName: repositoryName,
	}
}

func ConvertRepositoryResponseToRepository(response models.RepositoryResponse) models.Repository {
	description := ""
	if response.Description != nil {
		description = response.Description.(string)
	}

	return models.Repository{
		Name:            response.Name,
		Description:     description,
		URL:             response.HTMLURL,
		Language:        response.Language,
		ForksCount:      response.ForksCount,
		StarsCount:      response.StargazersCount,
		OpenIssuesCount: response.OpenIssuesCount,
		WatchersCount:   response.WatchersCount,
		CreatedAt:       response.CreatedAt,
		UpdatedAt:       response.UpdatedAt,
	}
}
