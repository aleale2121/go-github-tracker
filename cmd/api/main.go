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

const webPort = "8081"

func main() {
	
	githubRestClient := githubrestclient.NewGithubRestClient(&models.Config{
		GithubToken:    os.Getenv("GITHUB_TOKEN"),
		GithubUsername: os.Getenv("GITHUB_USERNAME"),
	})

	repositoriesHandler := handlers.NewRepositoriesHandler(githubRestClient)
	repositoriesRouting := routing.RepositoriesRouting(repositoriesHandler)

	var routesList []routers.Route
	routesList = append(routesList, repositoriesRouting...)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routers.Routes(routesList),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
