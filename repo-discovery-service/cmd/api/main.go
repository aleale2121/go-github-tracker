package main

import (
	"database/sql"
	"time"

	"fmt"
	"commits-manager-service/internal/constants/models"
	"commits-manager-service/internal/glue/routing"
	"commits-manager-service/internal/handlers"
	"commits-manager-service/internal/pkg/githubrestclient"
	"commits-manager-service/internal/storage/db"
	"commits-manager-service/platforms/routers"

	"commits-manager-service/internal/services/commitsmanagerservice"
	"commits-manager-service/internal/services/repomanagerservice"

	"commits-manager-service/internal/services/commitsmonitorservice"
	"commits-manager-service/internal/services/reposdiscoveryservice"

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
	repositoryManagerService := repomanagerservice.NewRepositoryManagerService(repositoryPersistence)
	repositoriesHandler := handlers.NewRepositoriesHandler(repositoryManagerService)
	repositoriesRouting := routing.RepositoriesRouting(repositoriesHandler)

	commitPersistence := db.NewCommitPersistence(dbConn)
	commitsManagerService := commitsmanagerservice.NewCommitsManagerService(commitPersistence)
	commitsHandler := handlers.NewCommitsHandler(commitsManagerService)
	commitsRouting := routing.CommitsRouting(commitsHandler)

	MetaDataPersistence := db.NewMetadataPersistence(dbConn)

	var routesList []routers.Route
	routesList = append(routesList, repositoriesRouting...)
	routesList = append(routesList, commitsRouting...)

	commitsMonitorService := commitsmonitorservice.NewCommentMonitorService(repositoryPersistence, commitPersistence,
		MetaDataPersistence, githubRestClient)

	reposDiscoveryService := reposdiscoveryservice.NewReposDiscoveryService(repositoryPersistence,
		MetaDataPersistence, githubRestClient)

	go reposDiscoveryService.ScheduleFetchingRepository(time.Hour * 24)
	go commitsMonitorService.ScheduleFetchingCommits(time.Hour * 1)

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