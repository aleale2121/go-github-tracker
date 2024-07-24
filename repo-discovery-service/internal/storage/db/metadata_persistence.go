package db

import (
	"database/sql"
	"go-github-tracker/internal/constants/models"
	"log"
	"time"
)

type MetadataPersistence struct {
	db *sql.DB
}

// NewMetadataPersistence creates an instance of the MetadataPersistence.
func NewMetadataPersistence(dbPool *sql.DB) MetadataPersistence {
	return MetadataPersistence{db: dbPool}
}

// SaveFetchReposMetadata saves metadata for fetching repositories.
func (mp *MetadataPersistence) SaveFetchReposMetadata(metadata models.FetchReposMetadata) error {
	stmt := `INSERT INTO fetch_repos_metadata (total, fetched_at) VALUES ($1, $2)`
	_, err := mp.db.Exec(stmt, metadata.Total, metadata.FetchedAt)
	if err != nil {
		log.Println("Error inserting fetch repos metadata:", err)
		return err
	}
	return nil
}

// SaveFetchCommitsMetadata saves metadata for fetching commits.
func (mp *MetadataPersistence) SaveFetchCommitsMetadata(metadata models.FetchCommitsMetadata) error {
	stmt := `INSERT INTO fetch_commits_metadata (repository_name, total, fetched_at) VALUES ($1, $2, $3)`
	_, err := mp.db.Exec(stmt, metadata.RepositoryName, metadata.Total, metadata.FetchedAt)
	if err != nil {
		log.Println("Error inserting fetch commits metadata:", err)
		return err
	}
	return nil
}

// GetLastReposFetchTime returns the last repository fetch time.
func (mp *MetadataPersistence) GetLastReposFetchTime() (time.Time, error) {
	var fetchedAt time.Time
	err := mp.db.QueryRow("SELECT fetched_at FROM fetch_repos_metadata ORDER BY fetched_at DESC LIMIT 1").Scan(&fetchedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil 
		}
		log.Println("Error querying last repository fetch time:", err)
		return time.Time{}, err
	}
	return fetchedAt, nil
}

// GetLastCommitFetchTime returns the last commit fetch time for a given repository.
func (mp *MetadataPersistence) GetLastCommitFetchTime(repositoryName string) (time.Time, error) {
	var fetchedAt time.Time
	err := mp.db.QueryRow("SELECT fetched_at FROM fetch_commits_metadata WHERE repository_name = $1 ORDER BY fetched_at DESC LIMIT 1", repositoryName).Scan(&fetchedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil 
		}
		log.Println("Error querying last commit fetch time:", err)
		return time.Time{}, err
	}
	return fetchedAt, nil
}
