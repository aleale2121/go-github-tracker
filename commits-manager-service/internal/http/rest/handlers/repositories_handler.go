package handlers

import (
	"errors"
	"fmt"
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

    totalRepositories, err := h.RepositoryPersistence.GetTotalRepositories()
    if err != nil {
        errorJSON(w, errors.New("failed to fetch total number of repositories"), http.StatusBadRequest)
        return
    }

    totalPages := (totalRepositories + limit - 1) / limit

    prevPage := ""
    if page > 1 {
        prevPage = fmt.Sprintf("/repositories?page=%d&limit=%d", page-1, limit)
    }

    nextPage := ""
    if page < totalPages {
        nextPage = fmt.Sprintf("/repositories?page=%d&limit=%d", page+1, limit)
    }

    payload := jsonResponse{
        Error:   false,
        Message: "repositories",
        Data:    repositories,
        Pagination: map[string]interface{}{
            "currentPage": page,
            "prevPage":    prevPage,
            "nextPage":    nextPage,
            "totalPages":  totalPages,
        },
    }

    writeJSON(w, http.StatusOK, payload)
}
