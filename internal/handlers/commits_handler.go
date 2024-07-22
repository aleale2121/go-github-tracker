package handlers

import (
	"errors"
	"net/http"

	"go-github-tracker/internal/storage/db"

	"github.com/go-chi/chi/v5"
)

type CommitsHandler struct {
	CommitPersistence db.CommitPersistence
}

func NewcommitsHandler(commitPersistence db.CommitPersistence) *CommitsHandler {
	return &CommitsHandler{
		CommitPersistence: commitPersistence,
	}
}

func (h *CommitsHandler) GetAllcommits(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repositoryName")

	commits, err := h.CommitPersistence.GetCommitsByRepoName(repoName)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch commits"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "commits",
		Data:    commits,
	}

	writeJSON(w, http.StatusAccepted, payload)
}
