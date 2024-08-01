package db_test

import (
	"commits-manager-service/internal/storage/db"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"testing"
)

var repositoryQueries db.GitReposRepository
var commitsQueries db.CommitRepository

func TestMain(m *testing.M) {

	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal("cannot open testDB:", err)
	}

	err = testDB.Ping()
	if err != nil {
		log.Fatal("cannot ping :", err)
	}

	createTablesQuery := `
	CREATE TABLE repositories
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) UNIQUE,
		description TEXT,
		url VARCHAR(255) NOT NULL,
		language VARCHAR(255),
		forks_count INT NOT NULL,
		stars_count INT NOT NULL,
		open_issues_count INT NOT NULL,
		watchers_count INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	
	CREATE TABLE commits
	(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sha VARCHAR(255) UNIQUE,
		url VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		author_name VARCHAR(255) NOT NULL,
		author_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		repository_name VARCHAR(255) NOT NULL,
		FOREIGN KEY (repository_name) REFERENCES repositories(name) ON DELETE CASCADE
	);`
	_, err = testDB.Exec(createTablesQuery)
	if err != nil {
		log.Fatal("Cannot create tables:", err)
	}

	repositoryQueries = db.NewRepositoryPersistence(testDB)
	commitsQueries = db.NewCommitPersistence(testDB)

	os.Exit(m.Run())
}
