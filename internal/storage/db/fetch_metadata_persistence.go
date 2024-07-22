package db

import (
	"database/sql"
	"go-github-tracker/internal/constants/models"
	"log"
)

type FetchMetadataPersistence struct {
	db *sql.DB
}

// NewFetchMetadataPersistence creates an instance of the FetchMetadataPersistence.
func NewFetchMetadataPersistence(dbPool *sql.DB) FetchMetadataPersistence {
	return FetchMetadataPersistence{db: dbPool}
}

// GetLastFetchMetadata returns the last fetch metadata.
func (fmp *FetchMetadataPersistence) GetLastFetchMetadata() (*models.FetchMetadata, error) {
	var metadata models.FetchMetadata
	err := fmp.db.QueryRow("SELECT id, last_repository_fetch, last_commit_fetch FROM fetch_metadata ORDER BY id DESC LIMIT 1").
		Scan(&metadata.ID, &metadata.LastRepositoryFetch, &metadata.LastCommitFetch)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, return default metadata
			return &models.FetchMetadata{}, nil
		}
		log.Println("Error querying fetch metadata:", err)
		return nil, err
	}
	return &metadata, nil
}

// UpdateLastFetchMetadata updates the last fetch metadata.
func (fmp *FetchMetadataPersistence) UpdateLastFetchMetadata(metadata models.FetchMetadata) error {
	_, err := fmp.db.Exec("INSERT INTO fetch_metadata (last_repository_fetch, last_commit_fetch) VALUES ($1, $2)",
		metadata.LastRepositoryFetch, metadata.LastCommitFetch)
	if err != nil {
		log.Println("Error updating fetch metadata:", err)
		return err
	}
	return nil
}
