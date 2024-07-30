package db

import (
	"context"
	"database/sql"
	"commits-manager-service/internal/constants/models"
	"log"
)

type GitReposRepository interface {
	GetAllRepositories() ([]*models.Repository, error)
	GetRepositoryByID(name string) (*models.Repository, error)
	UpdateRepository(repo models.Repository) error
	DeleteRepository(name string) error
	InsertRepository(repo models.Repository) (string, error)
	SaveAllRepositories(repos []models.Repository) error
	RepositoryExists(name string) (bool, error)
}
type RepositoryPersistence struct {
	db *sql.DB
}

// NewRepositoryPersistence creates an instance of the RepositoryPersistence.
func NewRepositoryPersistence(dbPool *sql.DB) GitReposRepository {
	return &RepositoryPersistence{db: dbPool}
}

// GetAllRepositories returns all repositories from the database.
func (rp *RepositoryPersistence) GetAllRepositories() ([]*models.Repository, error) {
	rows, err := rp.db.Query("SELECT * FROM repositories")
	if err != nil {
		log.Println("Error querying repositories:", err)
		return nil, err
	}
	defer rows.Close()

	var repositories []*models.Repository
	for rows.Next() {
		var repo models.Repository
		if err := rows.Scan(&repo.Name, &repo.Description, &repo.URL, &repo.Language, &repo.ForksCount, &repo.StarsCount, &repo.OpenIssuesCount, &repo.WatchersCount, &repo.CreatedAt, &repo.UpdatedAt); err != nil {
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
func (rp *RepositoryPersistence) GetRepositoryByID(name string) (*models.Repository, error) {
	var repo models.Repository
	err := rp.db.QueryRow("SELECT name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at FROM repositories WHERE name = $1", name).
		Scan(&repo.Name, &repo.Description, &repo.URL, &repo.Language, &repo.ForksCount, &repo.StarsCount, &repo.OpenIssuesCount, &repo.WatchersCount, &repo.CreatedAt, &repo.UpdatedAt)
	if err != nil {
		log.Println("Error querying repository by ID:", err)
		return nil, err
	}
	return &repo, nil
}

// UpdateRepository updates a repository in the database.
func (rp *RepositoryPersistence) UpdateRepository(repo models.Repository) error {
	_, err := rp.db.Exec("UPDATE repositories SET  description = $1, url = $2, language = $3, forks_count = $4, stars_count = $5, open_issues_count = $6, watchers_count = $7, created_at = $8, updated_at = $9 WHERE name = $10",
		repo.Description, repo.URL, repo.Language, repo.ForksCount, repo.StarsCount, repo.OpenIssuesCount, repo.WatchersCount, repo.CreatedAt, repo.UpdatedAt, repo.Name)
	if err != nil {
		log.Println("Error updating repository:", err)
		return err
	}
	return nil
}

// DeleteRepository deletes a repository from the database.
func (rp *RepositoryPersistence) DeleteRepository(name string) error {
	_, err := rp.db.Exec("DELETE FROM repositories WHERE name = $1", name)
	if err != nil {
		log.Println("Error deleting repository:", err)
		return err
	}
	return nil
}

// InsertRepository inserts a new repository into the database.
func (rp *RepositoryPersistence) InsertRepository(repo models.Repository) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO repositories (name, description, url, language, forks_count, stars_count, open_issues_count, watchers_count, created_at, updated_at) 
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning name`

	var name string
	err := rp.db.QueryRowContext(ctx, stmt, repo.Name, repo.Description, repo.URL, repo.Language, repo.ForksCount, repo.StarsCount, repo.OpenIssuesCount, repo.WatchersCount, repo.CreatedAt, repo.UpdatedAt).Scan(&name)
	if err != nil {
		log.Println("Error inserting repository:", err)
		return "", err
	}
	return name, nil
}

// SaveAllRepositories inserts or updates multiple repositories in the database.
func (rp *RepositoryPersistence) SaveAllRepositories(repos []models.Repository) error {
	for _, repo := range repos {
		exists, err := rp.RepositoryExists(repo.Name)
		if err != nil {
			return err
		}
		if exists {
			if err := rp.UpdateRepository(repo); err != nil {
				return err
			}
		} else {
			if _, err := rp.InsertRepository(repo); err != nil {
				return err
			}
		}
	}
	return nil
}

// RepositoryExists checks if a repository exists in the database.
func (rp *RepositoryPersistence) RepositoryExists(name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM repositories WHERE name = $1)"
	err := rp.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}