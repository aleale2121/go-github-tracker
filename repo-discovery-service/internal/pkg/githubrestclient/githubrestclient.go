package githubrestclient

import (
	"encoding/json"
	"fmt"
	"repos-discovery-service/internal/constants/models"
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

func (gp GithubRestClient) FetchRepositories(since string) ([]models.RepositoryResponse, error) {
	fetchRepoUrl := base_url + fmt.Sprintf("/users/%s/repos?sort=updated&direction=desc", gp.Config.GithubUsername)
	if since != "" {
		fetchRepoUrl += fmt.Sprintf("&since=%s", since)
	}
	fmt.Println(fetchRepoUrl)
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
	log.Println(response.StatusCode)
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil, err
	}

	// Print the response body as a string
	// bodyString := string(bodyBytes)
	// fmt.Println("Fetch Response Body: ", bodyString)

	var repositories []models.RepositoryResponse
	err = json.Unmarshal(bodyBytes, &repositories)
	if err != nil {
		log.Println("Error unmarshalling response body:", err)
		return nil, err
	}

	return repositories, nil
}
