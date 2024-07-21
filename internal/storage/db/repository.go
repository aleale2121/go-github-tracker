package db

import (
	"context"
	"database/sql"
	"go-github-tracker/internal/constants/models"
	"log"
)

type RepositoryPersistence struct {
	db *sql.DB
}

// NewRepositoryPersistence creates an instance of the RepositoryPersistence.
func NewRepositoryPersistence(dbPool *sql.DB) RepositoryPersistence {
	return RepositoryPersistence{db: dbPool}
}

// GetAllRepositories returns all repositories from the database.
func (rp *RepositoryPersistence) GetAllRepositories() ([]*models.Repository, error) {
	rows, err := rp.db.Query("SELECT id, name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at FROM repositories")
	if err != nil {
		log.Println("Error querying repositories:", err)
		return nil, err
	}
	defer rows.Close()

	var repositories []*models.Repository
	for rows.Next() {
		var repo models.Repository
		if err := rows.Scan(&repo.ID, &repo.Name, &repo.Description, &repo.URL, &repo.Language, &repo.ForksCount, &repo.StarsCount, &repo.OpenIssuesCount, &repo.WatchersCount, &repo.CreatedAt, &repo.UpdatedAt); err != nil {
			log.Println("Error scanning repository row:", err)
			return nil, err
		}
		repositories = append(repositories, &repo)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through repositories:", err)
		return nil, err
	}

	return repositories, nil
}

// GetRepositoryByID returns a repository from the database by ID.
func (rp *RepositoryPersistence) GetRepositoryByID(id int64) (*models.Repository, error) {
	var repo models.Repository
	err := rp.db.QueryRow("SELECT id, name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at FROM repositories WHERE id = $1", id).
		Scan(&repo.ID, &repo.Name, &repo.Description, &repo.URL, &repo.Language, &repo.ForksCount, &repo.StarsCount, &repo.OpenIssuesCount, &repo.WatchersCount, &repo.CreatedAt, &repo.UpdatedAt)
	if err != nil {
		log.Println("Error querying repository by ID:", err)
		return nil, err
	}
	return &repo, nil
}

// UpdateRepository updates a repository in the database.
func (rp *RepositoryPersistence) UpdateRepository(repo models.Repository) error {
	_, err := rp.db.Exec("UPDATE repositories SET name = $1, description = $2, url = $3, language = $4, forks_count = $5, stars_count = $6, open_issues_count = $7, watchers_count = $8, created_at = $9, updated_at = $10 WHERE id = $11",
		repo.Name, repo.Description, repo.URL, repo.Language, repo.ForksCount, repo.StarsCount, repo.OpenIssuesCount, repo.WatchersCount, repo.CreatedAt, repo.UpdatedAt, repo.ID)
	if err != nil {
		log.Println("Error updating repository:", err)
		return err
	}
	return nil
}

// DeleteRepository deletes a repository from the database.
func (rp *RepositoryPersistence) DeleteRepository(id int64) error {
	_, err := rp.db.Exec("DELETE FROM repositories WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting repository:", err)
		return err
	}
	return nil
}

// InsertRepository inserts a new repository into the database.
func (rp *RepositoryPersistence) InsertRepository(repo models.Repository) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO repositories (name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at) 
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	var id int64
	err := rp.db.QueryRowContext(ctx, stmt, repo.Name, repo.Description, repo.URL, repo.Language, repo.ForksCount, repo.StarsCount, repo.OpenIssuesCount, repo.WatchersCount, repo.CreatedAt, repo.UpdatedAt).Scan(&id)
	if err != nil {
		log.Println("Error inserting repository:", err)
		return 0, err
	}
	return id, nil
}

// Function to check if a repository exists
func (rp *RepositoryPersistence) RepositoryExists(db *sql.DB, id int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM repositories WHERE id = $1)"
	err := rp.db.QueryRow(query, id).Scan(&exists)
	return exists, err
}
