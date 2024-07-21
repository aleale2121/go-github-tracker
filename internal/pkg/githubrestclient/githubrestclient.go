package githubrestclient

import (
	"encoding/json"
	"fmt"
	"go-github-tracker/internal/constants/models"
	"log"
	"time"

	"net/http"
)

type GithubRestClient struct {
	Config *models.Config
}

func NewGithubRestClient(Config *models.Config) GithubRestClient {
	return GithubRestClient{Config: Config}
}

const base_url = "https://api.github.com"

func (gp GithubRestClient) FetchRepositories() ([]models.RepositoryReponse, error) {
	fetchRepoUrl := base_url + fmt.Sprintf("/users/%s/repos", gp.Config.GithubUsername)

	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	var repositories []models.RepositoryReponse
	err = json.NewDecoder(response.Body).Decode(&repositories)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return repositories, nil
}

func (gp GithubRestClient) FetchCommits(repositoryName string, since time.Time) ([]models.CommitReponse, error) {
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


