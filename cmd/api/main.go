package main

import (
	"database/sql"
	"time"

	"fmt"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/glue/routing"
	"go-github-tracker/internal/handlers"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/internal/pkg/scheduler"
	"go-github-tracker/internal/storage/db"
	"go-github-tracker/platforms/routers"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

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

	githubRestClient := githubrestclient.NewGithubRestClient(&models.Config{
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		GithubUsername: os.Getenv("GITHUB_USERNAME"),
	})

	repositoryPersistence := db.NewRepositoryPersistence(dbConn)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoryPersistence)
	repositoriesRouting := routing.RepositoriesRouting(repositoriesHandler)

	commitPersistence := db.NewCommitPersistence(dbConn)
	commitsHandler := handlers.NewcommitsHandler(commitPersistence)
	commitsRouting := routing.CommitsRouting(commitsHandler)

	MetaDataPersistence := db.NewMetadataPersistence(dbConn)

	var routesList []routers.Route
	routesList = append(routesList, repositoriesRouting...)
	routesList = append(routesList, commitsRouting...)

	schedulerService := scheduler.NewSchedulerService(repositoryPersistence, commitPersistence,
		MetaDataPersistence, githubRestClient)

	go schedulerService.ScheduleFetchingRepository(time.Second * 30)
	go schedulerService.ScheduleFetchingCommits(time.Second * 60)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routers.Routes(routesList),
	}
	log.Println("server started at port :80")

	err := srv.ListenAndServe()
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
