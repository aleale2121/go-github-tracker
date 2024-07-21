package handlers

import (
	"errors"
	"net/http"
	"time"

	"go-github-tracker/internal/pkg/githubrestclient"

	"github.com/go-chi/chi/v5"
)

type CommitsHandler struct {
	GithubRestClient githubrestclient.GithubRestClient
}

func NewcommitsHandler(githubrestclient githubrestclient.GithubRestClient) *CommitsHandler {
	return &CommitsHandler{
		GithubRestClient: githubrestclient,
	}
}

func (h *CommitsHandler) GetAllcommits(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repositoryName")

	commits, err := h.GithubRestClient.FetchCommits(repoName, time.Time{})
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
