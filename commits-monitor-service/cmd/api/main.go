package main

import (
	"commits-monitor-service/internal/constants"
	"commits-monitor-service/internal/constants/models"
	"commits-monitor-service/internal/http/grpc/client/commits"
	"commits-monitor-service/internal/http/grpc/client/repos"
	"commits-monitor-service/internal/pkg/githubrestclient"
	"fmt"
	"math"
	"time"

	"commits-monitor-service/internal/services/commitsmonitorservice"

	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

const commitMangerUrl = "commits-manager-service:50001"

func main() {
	startDate := "1970-10-03T10:01:20Z"
	_, err := time.Parse(constants.ISO_8601_TIME_LAYOUT, os.Getenv("START_DATE"))
	if err != nil {
		log.Println("Cannot parse start date: ", err)
	} else {
		startDate = os.Getenv("START_DATE")
	}

	endDate := "2098-10-03T10:01:20Z"
	_, err = time.Parse(constants.ISO_8601_TIME_LAYOUT, os.Getenv("END_DATE"))
	if err != nil {
		log.Println("Cannot parse end date: ", err)
	} else {
		endDate = os.Getenv("END_DATE")
	}
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
		StartDate:      startDate,
		EndDate:        endDate,
	})

	commitMetaDataServiceClient := commits.NewCommitsMetaDataServiceClient(commitMangerUrl)
	reposMetaDataServiceClient := repos.NewReposMetaDataServiceClient(commitMangerUrl)
	commitsMonitorService := commitsmonitorservice.NewCommentMonitorService(githubRestClient,
		*reposMetaDataServiceClient, *commitMetaDataServiceClient, rabbitConn)

	wait := make(chan bool)

	// Wait 1 minute for first repository fetch
	timer := time.After(60 * time.Second)
	<-timer

	go commitsMonitorService.ScheduleFetchingCommits(time.Hour * 1)

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
