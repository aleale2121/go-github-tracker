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
	rows, err := cp.db.Query("SELECT sha, url, message, author_name, author_date FROM commits")
	if err != nil {
		log.Println("Error querying commits:", err)
		return nil, err
	}
	defer rows.Close()

	var commits []*models.Commit
	for rows.Next() {
		var commit models.Commit
		if err := rows.Scan(&commit.SHA, &commit.URL, &commit.Message, &commit.AuthorName, &commit.AuthorDate); err != nil {
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
	err := cp.db.QueryRow("SELECT sha, url, message, author_name, author_date FROM commits WHERE sha = $1", sha).
		Scan(&commit.SHA, &commit.URL, &commit.Message, &commit.AuthorName, &commit.AuthorDate)
	if err != nil {
		log.Println("Error querying commit by SHA:", err)
		return nil, err
	}
	return &commit, nil
}

// UpdateCommit updates a commit in the database.
func (cp *CommitPersistence) UpdateCommit(commit models.Commit) error {
	_, err := cp.db.Exec("UPDATE commits SET url = $1, message = $2, author_name = $3, author_date = $4 WHERE sha = $5",
		commit.URL, commit.Message, commit.AuthorName, commit.AuthorDate, commit.SHA)
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
func (cp *CommitPersistence) InsertCommit(commit models.Commit) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO commits (sha, url, message, author_name, author_date) 
             VALUES ($1, $2, $3, $4, $5) returning sha`

	var sha string
	err := cp.db.QueryRowContext(ctx, stmt, commit.SHA, commit.URL, commit.Message, commit.AuthorName, commit.AuthorDate).Scan(&sha)
	if err != nil {
		log.Println("Error inserting commit:", err)
		return "", err
	}
	return sha, nil
}

// Function to check if a commit exists
func (cp *CommitPersistence) CommitExists(db *sql.DB, sha string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM commits WHERE sha = $1)"
	err := cp.db.QueryRow(query, sha).Scan(&exists)
	return exists, err
}
