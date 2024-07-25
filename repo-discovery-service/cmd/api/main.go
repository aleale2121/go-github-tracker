package main

import (
	"repos-discovery-service/internal/constants/models"
	"repos-discovery-service/internal/http/grpc/client/repos"
	"repos-discovery-service/internal/pkg/githubrestclient"

	"fmt"
	"math"
	"time"

	"repos-discovery-service/internal/services/reposdiscoveryservice"

	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const commitMangerUrl = "commits-manager-service:50001"

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	githubRestClient := githubrestclient.NewGithubRestClient(&models.Config{
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		GithubUsername: os.Getenv("GITHUB_USERNAME"),
	})

	reposMetaDataServiceClient := repos.NewReposMetaDataServiceClient(commitMangerUrl)
	reposdiscoveryservice := reposdiscoveryservice.NewReposDiscoveryService(githubRestClient,
		*reposMetaDataServiceClient,
		rabbitConn)

	wait := make(chan bool)

	// Wait one minute until commit-manager service started
    timer := time.After(30 * time.Second)
	<-timer

	go reposdiscoveryservice.ScheduleFetchingRepository(time.Hour * 1)

	<-wait

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
