package handlers

import (
	"errors"
	"go-github-tracker/internal/pkg/githubapi"
	"net/http"
)

type RepositoriesHandler struct {
	GithubApi githubapi.GithubAPi
}

func NewRepositoriesHandler(githubapi githubapi.GithubAPi) *RepositoriesHandler {
	return &RepositoriesHandler{
		GithubApi: githubapi,
	}
}

func (h *RepositoriesHandler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := h.GithubApi.FetchRepositories()
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
