package githubrestclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"commits-monitor-service/internal/constants/models"
)

type GithubRestClient struct {
	Config *models.Config
}

func NewGithubRestClient(Config *models.Config) GithubRestClient {
	return GithubRestClient{Config: Config}
}

const baseURL = "https://api.github.com"

// buildURI constructs the URL with query parameters.
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

func (gp GithubRestClient) FetchCommits(repositoryName, since string, perPage, page int) ([]models.CommitResponse, error) {
	path := fmt.Sprintf("/repos/%s/%s/commits", gp.Config.GithubUsername, repositoryName)
	queryParams := map[string]string{}
	if since != "" {
		queryParams["since"] = since
	}
	queryParams["per_page"] = fmt.Sprintf("%d", perPage)
	queryParams["page"] = fmt.Sprintf("%d", page)

	fetchRepoUrl := buildURI(baseURL, path, queryParams)

	request, err := http.NewRequest(http.MethodGet, fetchRepoUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gp.Config.GithubToken))
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("CMOS: Error reading response body:", err)
		return nil, err
	}

	var commits []models.CommitResponse
	err = json.Unmarshal(bodyBytes, &commits)
	if err != nil {
		log.Println("CMOS: Error unmarshalling response body:", err)
		return nil, err
	}

	return commits, nil
}
