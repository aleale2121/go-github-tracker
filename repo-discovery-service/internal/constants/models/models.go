package models

type Config struct {
	DSN            string `json:"dsn"`
	GithubToken    string `json:"github_token"`
	GithubUsername string `json:"github_username"`
}

