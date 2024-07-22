# GitHub API Data Fetching Service

## Project Overview

This service is built in Go and interacts with GitHub's public APIs to fetch repository and commit data. It stores this data in a PostgreSQL database and continuously monitors for changes.

## Implementation Details

### Fetching GitHub API Data

- **Commits:**
  - Fetches commit message, author, date, and URL.
  - Saves data in PostgreSQL to ensure commits in the database mirror those on GitHub.
  - Uses a configurable date to start fetching commits and allows resetting the collection.

- **Repository Metadata:**
  - Stores repository details such as name, description, URL, language, forks count, stars count, open issues count, watchers count, and creation/update dates.

### Data Storage

- Data is stored in PostgreSQL with tables for repositories, commits, and metadata.
- Efficient querying ensures performance and scalability.

### Scheduling

- The service uses Go's `time.Ticker` to schedule data fetching at regular intervals (e.g., every hour).
- Prevents duplicate commits by comparing fetched data with existing records in the database.

### Endpoints

- **List Repositories:**
    GET <http://localhost:8081/repositories>
    Retrieves a list of all repositories.

- **Fetch Repository Commits:**
    GET <http://localhost:8081/commits/{repoName}>
    Retrieves commits for a specific repository.

- **Fetch Overall Top N Committers:**
    GET <http://localhost:8081/top-commit-authors?limit=10>
    Retrieves the top N commit authors overall.
- **Fetch Top N Committers for a Specific Repository:**
    GET <http://localhost:8081/top-commit-authors/{repoName}?limit=10>  
    Retrieves the top N commit authors for a specific repository.

### Unit Tests

- Unit tests are included to validate core functionalities, such as data fetching and persistence.

## Setup and Usage

1. **Update Environment Variables:**

    - Update `app.env` with your GitHub token and username:

    ```markdown
    DSN=host=postgres port=5432 user=postgres password=password dbname=github_tracker sslmode=disable timezone=UTC connect_timeout=5
    GITHUB_TOKEN=""
    GITHUB_USERNAME=""
    ```

2. **Build and Run:**

    - Use Docker to build and start the services:

    ```markdown
    docker-compose up --build
    ```

Ensure to review the code, update the environment variables, and follow the provided instructions for a successful setup and execution of the service.
