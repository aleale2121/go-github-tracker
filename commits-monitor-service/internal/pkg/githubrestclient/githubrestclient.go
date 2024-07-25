package githubrestclient

import (
	"encoding/json"
	"fmt"
	"commits-monitor-service/internal/constants/models"
	"log"
	"io"

	"net/http"
)

type GithubRestClient struct {
	Config *models.Config
}

func NewGithubRestClient(Config *models.Config) GithubRestClient {
	return GithubRestClient{Config: Config}
}

const base_url = "https://api.github.com"

func (gp GithubRestClient) FetchCommits(repositoryName, since string) ([]models.CommitResponse, error) {
	fetchRepoUrl := base_url + fmt.Sprintf("/repos/%s/%s/commits", gp.Config.GithubUsername, repositoryName)
	if since != "" {
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

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("cm Error reading response body:", err)
		return nil, err
	}

	fmt.Println("cm Response Status Code:", response.StatusCode)

	var repositories []models.CommitResponse
	err = json.Unmarshal(bodyBytes, &repositories)
	if err != nil {
		log.Println("cm Error unmarshalling response body:", err)
		return nil, err
	}

	return repositories, nil
}
