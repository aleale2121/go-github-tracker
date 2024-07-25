package handlers

import (
	"errors"
	"net/http"

	"commits-manager-service/internal/module/repos"

)

type RepositoriesHandler struct {
	RepositoryPersistence repos.RepositoryManagerService
	
}

func NewRepositoriesHandler(repositoryPersistence repos.RepositoryManagerService) *RepositoriesHandler {
	return &RepositoriesHandler{
		RepositoryPersistence: repositoryPersistence,
	}
}

func (h *RepositoriesHandler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := h.RepositoryPersistence.GetRepositories()
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
