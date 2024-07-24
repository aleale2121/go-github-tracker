package main

import (
	"database/sql"
	"math"
	"time"

	"commits-manager-service/internal/glue/routing"
	"commits-manager-service/internal/http/rest/handlers"
	event "commits-manager-service/internal/message-broker/rabbitmq"
	"commits-manager-service/internal/storage/db"
	"commits-manager-service/platforms/routers"
	"fmt"

	"commits-manager-service/internal/services/commitsmanagerservice"
	"commits-manager-service/internal/services/repomanagerservice"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	amqp "github.com/rabbitmq/amqp091-go"

	"log"
	"net/http"
	"os"
)

const webPort = "80"

var counts int64

func main() {
	dbConn := connectToDB()
	if dbConn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	repositoryPersistence := db.NewRepositoryPersistence(dbConn)
	repositoryManagerService := repomanagerservice.NewRepositoryManagerService(repositoryPersistence)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoryManagerService)
	repositoriesRouting := routing.RepositoriesRouting(repositoriesHandler)

	commitPersistence := db.NewCommitPersistence(dbConn)
	commitsManagerService := commitsmanagerservice.NewCommitsManagerService(commitPersistence)
	commitsHandler := handlers.NewCommitsHandler(commitsManagerService)
	commitsRouting := routing.CommitsRouting(commitsHandler)

	var routesList []routers.Route
	routesList = append(routesList, repositoriesRouting...)
	routesList = append(routesList, commitsRouting...)

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Println("Listening for and consuming RabbitMQ messages...")
		panic(err)
	}

	// watch the queue and consume events
	go func(eventConsumer event.Consumer) {
		err = eventConsumer.Listen([]string{"github.REPOS", "github.COMMITS"})
		if err != nil {
			log.Println(err)
		}
	}(consumer)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routers.Routes(routesList),
	}
	log.Println("server started at port :80")

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	fmt.Println("DSN-->", dsn)
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
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
