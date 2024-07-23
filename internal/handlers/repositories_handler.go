package handlers

import (
	"errors"
	"net/http"

	"go-github-tracker/internal/services/repomanagerservice"

)

type RepositoriesHandler struct {
	RepositoryPersistence repomanagerservice.RepositoryManagerService
	
}

func NewRepositoriesHandler(repositoryPersistence repomanagerservice.RepositoryManagerService) *RepositoriesHandler {
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
