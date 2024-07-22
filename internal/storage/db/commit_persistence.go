package db

import (
	"context"
	"database/sql"
	"go-github-tracker/internal/constants/models"
	"log"
)

type CommitPersistence struct {
	db *sql.DB
}

// NewCommitPersistence creates an instance of the CommitPersistence.
func NewCommitPersistence(dbPool *sql.DB) CommitPersistence {
	return CommitPersistence{db: dbPool}
}

// GetAllCommits returns all commits from the database.
func (cp *CommitPersistence) GetAllCommits() ([]*models.Commit, error) {
	rows, err := cp.db.Query("SELECT sha, url, message, author_name, author_date, created_at, updated_at, repository_name FROM commits")
	if err != nil {
		log.Println("Error querying commits:", err)
		return nil, err
	}
	defer rows.Close()

	var commits []*models.Commit
	for rows.Next() {
		var commit models.Commit
		if err := rows.Scan(&commit.SHA, &commit.URL, &commit.Message, &commit.AuthorName, &commit.AuthorDate, &commit.CreatedAt, &commit.UpdatedAt, &commit.RepositoryName); err != nil {
			log.Println("Error scanning commit row:", err)
			return nil, err
		}
		commits = append(commits, &commit)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through commits:", err)
		return nil, err
	}

	return commits, nil
}

// GetCommitBySHA returns a commit from the database by SHA.
func (cp *CommitPersistence) GetCommitBySHA(sha string) (*models.Commit, error) {
	var commit models.Commit
	err := cp.db.QueryRow("SELECT sha, url, message, author_name, author_date, created_at, updated_at, repository_name FROM commits WHERE sha = $1", sha).
		Scan(&commit.SHA, &commit.URL, &commit.Message, &commit.AuthorName, &commit.AuthorDate, &commit.CreatedAt, &commit.UpdatedAt, &commit.RepositoryName)
	if err != nil {
		log.Println("Error querying commit by SHA:", err)
		return nil, err
	}
	return &commit, nil
}

// UpdateCommit updates a commit in the database.
func (cp *CommitPersistence) UpdateCommit(commit models.Commit) error {
	_, err := cp.db.Exec("UPDATE commits SET url = $1, message = $2, author_name = $3, author_date = $4, created_at = $5, updated_at = $6, repository_name = $7 WHERE sha = $8",
		commit.URL, commit.Message, commit.AuthorName, commit.AuthorDate, commit.CreatedAt, commit.UpdatedAt, commit.RepositoryName, commit.SHA)
	if err != nil {
		log.Println("Error updating commit:", err)
		return err
	}
	return nil
}

// DeleteCommit deletes a commit from the database.
func (cp *CommitPersistence) DeleteCommit(sha string) error {
	_, err := cp.db.Exec("DELETE FROM commits WHERE sha = $1", sha)
	if err != nil {
		log.Println("Error deleting commit:", err)
		return err
	}
	return nil
}

// InsertCommit inserts a new commit into the database.
func (cp *CommitPersistence) InsertCommit(commit models.Commit) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO commits (sha, url, message, author_name, author_date, created_at, updated_at, repository_name) 
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := cp.db.ExecContext(ctx, stmt, commit.SHA, commit.URL, commit.Message, commit.AuthorName, commit.AuthorDate, commit.CreatedAt, commit.UpdatedAt, commit.RepositoryName)
	if err != nil {
		log.Println("Error inserting commit:", err)
		return err
	}
	return nil
}

// SaveAllCommits inserts or updates multiple commits in the database.
func (cp *CommitPersistence) SaveAllCommits(commits []models.Commit) error {
	for _, commit := range commits {
		exists, err := cp.CommitExists(commit.SHA)
		if err != nil {
			return err
		}
		if exists {
			if err := cp.UpdateCommit(commit); err != nil {
				return err
			}
		} else {
			if err := cp.InsertCommit(commit); err != nil {
				return err
			}
		}
	}
	return nil
}

// CommitExists checks if a commit exists in the database.
func (cp *CommitPersistence) CommitExists(sha string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM commits WHERE sha = $1)"
	err := cp.db.QueryRow(query, sha).Scan(&exists)
	return exists, err
}

// GetCommitsByRepoName returns all commits for a specific repository.
func (cp *CommitPersistence) GetCommitsByRepoName(repoName string) ([]*models.Commit, error) {
	rows, err := cp.db.Query("SELECT sha, url, message, author_name, author_date, created_at, updated_at, repository_name FROM commits WHERE repository_name = $1", repoName)
	if err != nil {
		log.Println("Error querying commits by repository name:", err)
		return nil, err
	}
	defer rows.Close()

	var commits []*models.Commit
	for rows.Next() {
		var commit models.Commit
		if err := rows.Scan(&commit.SHA, &commit.URL, &commit.Message, &commit.AuthorName, &commit.AuthorDate, &commit.CreatedAt, &commit.UpdatedAt, &commit.RepositoryName); err != nil {
			log.Println("Error scanning commit row:", err)
			return nil, err
		}
		commits = append(commits, &commit)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating through commits:", err)
		return nil, err
	}

	return commits, nil
}
