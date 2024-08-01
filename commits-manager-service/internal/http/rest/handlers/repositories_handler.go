package handlers

import (
	"errors"
	"net/http"
	"strconv"

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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	
	repositories, err := h.RepositoryPersistence.GetRepositories(limit, offset)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch repositories"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "repositories",
		Data:    repositories,
	}

	writeJSON(w, http.StatusOK, payload)
}
