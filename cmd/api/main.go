package main

import (
	"fmt"
	"go-github-tracker/internal/constants/models"
	"go-github-tracker/internal/glue/routing"
	"go-github-tracker/internal/handlers"
	"go-github-tracker/internal/pkg/githubrestclient"
	"go-github-tracker/platforms/routers"
	"log"
	"net/http"
	"os"
)

const webPort = "80"

func main() {
	fmt.Println("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	githubRestClient := githubrestclient.NewGithubRestClient(&models.Config{
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		GithubUsername: os.Getenv("GITHUB_USERNAME"),
	})

	repositoriesHandler := handlers.NewRepositoriesHandler(githubRestClient)
	repositoriesRouting := routing.RepositoriesRouting(repositoriesHandler)

	commitsHandler := handlers.NewcommitsHandler(githubRestClient)
	commitsRouting := routing.CommitsRouting(commitsHandler)

	var routesList []routers.Route
	routesList = append(routesList, repositoriesRouting...)
	routesList = append(routesList, commitsRouting...)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routers.Routes(routesList),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

	log.Println("server started at port :80")
}
