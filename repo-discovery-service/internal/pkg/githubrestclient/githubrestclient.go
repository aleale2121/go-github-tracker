package githubrestclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"repos-discovery-service/internal/constants/models"
)

type GithubRestClient struct {
	Config *models.Config
}

func NewGithubRestClient(Config *models.Config) GithubRestClient {
	return GithubRestClient{Config: Config}
}

const baseURL = "https://api.github.com"

func buildURI(base string, path string, queryParams map[string]string) string {
	u, _ := url.Parse(base)
	u.Path = path
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (gp GithubRestClient) FetchRepositories(perPage, page int) ([]models.RepositoryResponse, error) {
	path := fmt.Sprintf("/users/%s/repos", gp.Config.GithubUsername)
	queryParams := map[string]string{
		"sort":      "created",
		"direction": "desc",
		"per_page":  fmt.Sprintf("%d", perPage),
		"page":      fmt.Sprintf("%d", page),
	}

	fetchRepoUrl := buildURI(baseURL, path, queryParams)

	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	if err != nil {
		log.Println("RDS: ", err)
		return nil, err
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("RDS: ", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("RDS: unexpected status code: ", response.StatusCode)
		return nil, fmt.Errorf("RDS: unexpected status code: %d", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("RDS: error reading response body: ", err)
		return nil, err
	}

	var repositories []models.RepositoryResponse
	err = json.Unmarshal(bodyBytes, &repositories)
	if err != nil {
		log.Println("RDS: error unmarshalling response body: ", err)
		return nil, err
	}

	return repositories, nil
}

func (gp GithubRestClient) FetchRepositoryMetadata(repoName string) (models.RepositoryResponse, error) {
	path := fmt.Sprintf("/repos/%s/%s", gp.Config.GithubUsername, repoName)

	fetchRepoUrl := buildURI(baseURL, path, nil)

	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	if err != nil {
		log.Println("RDS: ", err)
		return models.RepositoryResponse{}, err
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("RDS: ", err)
		return models.RepositoryResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("RDS: fetch metaData unexpected status code: ", response.StatusCode)
		return models.RepositoryResponse{}, fmt.Errorf("RDS: unexpected status code: %d", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("RDS: error reading response body: ", err)
		return models.RepositoryResponse{}, err
	}

	var repository models.RepositoryResponse
	err = json.Unmarshal(bodyBytes, &repository)
	if err != nil {
		log.Println("RDS: error unmarshalling response body: ", err)
		return models.RepositoryResponse{}, err
	}

	return repository, nil
}
