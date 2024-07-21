package handlers

import (
	"errors"
	"go-github-tracker/internal/pkg/githubrestclient"
	"net/http"
)

type RepositoriesHandler struct {
	GithubRestClient githubrestclient.GithubRestClient
}

func NewRepositoriesHandler(githubrestclient githubrestclient.GithubRestClient) *RepositoriesHandler {
	return &RepositoriesHandler{
		GithubRestClient: githubrestclient,
	}
}

func (h *RepositoriesHandler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := h.GithubRestClient.FetchRepositories()
	if err != nil {
		errorJSON(w, errors.New("failed to fetch repositories"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "repositories",
		Data:    repositories,
	}

	writeJSON(w, http.StatusAccepted, payload)
}
