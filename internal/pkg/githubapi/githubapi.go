package githubapi

import (
	"encoding/json"
	"fmt"
	"go-github-tracker/internal/constants/models"
	"time"

	"net/http"
)

type GithubAPi struct {
	Config *models.Config
}

func NewGithubAPi(Config *models.Config) GithubAPi {
	return GithubAPi{Config: Config}
}

const base_url = "https://api.github.com"

func (gp GithubAPi) FetchRepositories() ([]models.RepositoryReponse, error) {
	fetchRepoUrl := base_url + fmt.Sprintf("/users/%s/repos", gp.Config.GithubUsername)

	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var repositories []models.RepositoryReponse
	err = json.NewDecoder(response.Body).Decode(&repositories)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func (gp GithubAPi) FetchCommits(repositoryName string, since time.Time) ([]models.CommitReponse, error) {
	fetchRepoUrl := base_url + fmt.Sprintf("/repos/%s/%s/commits", gp.Config.GithubUsername, repositoryName)
	if !since.IsZero() {
		fetchRepoUrl += fmt.Sprintf("?since=%s", since)
	}
	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var commits []models.CommitReponse
	err = json.NewDecoder(response.Body).Decode(&commits)
	if err != nil {
		return nil, err
	}

	return commits, nil
}


