package handlers

import (
	"errors"
	"go-github-tracker/internal/storage/db"
	"net/http"
)

type RepositoriesHandler struct {
	RepositoryPersistence db.RepositoryPersistence
	
}

func NewRepositoriesHandler(repositoryPersistence db.RepositoryPersistence) *RepositoriesHandler {
	return &RepositoriesHandler{
		RepositoryPersistence: repositoryPersistence,
	}
}

func (h *RepositoriesHandler) GetAllRepositories(w http.ResponseWriter, r *http.Request) {
	repositories, err := h.RepositoryPersistence.GetAllRepositories()
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
