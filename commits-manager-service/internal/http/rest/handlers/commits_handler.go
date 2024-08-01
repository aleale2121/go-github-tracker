package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"commits-manager-service/internal/module/commits"

	"github.com/go-chi/chi/v5"
)

type CommitsHandler struct {
	CommitsManagerService commits.CommitsManagerService
}

func NewCommitsHandler(commitPersistence commits.CommitsManagerService) *CommitsHandler {
	return &CommitsHandler{
		CommitsManagerService: commitPersistence,
	}
}

func (h *CommitsHandler) GetAllCommits(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repositoryName")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			errorJSON(w, errors.New("invalid startDate format"), http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Time{} 
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			errorJSON(w, errors.New("invalid endDate format"), http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Now() 
	}

	commits, err := h.CommitsManagerService.GetCommitsByRepositoryName(repoName, limit, offset, startDate, endDate)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch commits"), http.StatusBadRequest)
		return
	}

	totalCommits, err := h.CommitsManagerService.GetTotalCommitsByRepositoryName(repoName, startDate, endDate)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch total number of commits"), http.StatusBadRequest)
		return
	}

	totalPages := (totalCommits + limit - 1) / limit 

	prevPage := ""
	if page > 1 {
		prevPage = fmt.Sprintf("/repositories/%s/commits?page=%d&limit=%d&startDate=%s&endDate=%s", repoName, page-1, limit, startDateStr, endDateStr)
	}

	nextPage := ""
	if page < totalPages {
		nextPage = fmt.Sprintf("/repositories/%s/commits?page=%d&limit=%d&startDate=%s&endDate=%s", repoName, page+1, limit, startDateStr, endDateStr)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "commits",
		Data:    commits,
		Pagination: map[string]interface{}{
			"currentPage": page,
			"prevPage":    prevPage,
			"nextPage":    nextPage,
			"totalPages":  totalPages,
		},
	}

	writeJSON(w, http.StatusOK, payload)
}

func (h *CommitsHandler) GetTopCommitAuthors(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorJSON(w, errors.New("invalid limit parameter"), http.StatusBadRequest)
		return
	}

	authors, err := h.CommitsManagerService.GetTopCommitAuthors(limit)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch top commit authors"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "top commit authors",
		Data:    authors,
	}

	writeJSON(w, http.StatusOK, payload)
}

func (h *CommitsHandler) GetTopCommitAuthorsByRepo(w http.ResponseWriter, r *http.Request) {
	repoName := chi.URLParam(r, "repositoryName")
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errorJSON(w, errors.New("invalid limit parameter"), http.StatusBadRequest)
		return
	}

	authors, err := h.CommitsManagerService.GetTopCommitAuthorsByRepoName(repoName, limit)
	if err != nil {
		errorJSON(w, errors.New("failed to fetch top commit authors for repository"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "top commit authors for repository",
		Data:    authors,
	}

	writeJSON(w, http.StatusOK, payload)
}
