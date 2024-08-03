# GitHub API Data Fetching Service

## Project Overview

This project comprises three microservices designed to fetch, store, and monitor GitHub repository and commit data. The microservices interact with each other via gRPC and RabbitMQ to maintain a consistent and up-to-date database.

## Microservices Overview

### 1. Commits Manager Service

- **Functionality**:
  - Provides REST APIs and database implementations for repository and commit metadata.
  - Uses RabbitMQ to listen for new repository, new commit fetched events and metadata updates.
  - Runs a gRPC server to handle metadata queries from other microservices.
  
- **Responsibilities**:
  - Maintain the primary database.
  - Process events related to repository and commit updates.
  - Serve metadata to other services via gRPC.

### 2. Repos Discovery Service

- **Functionality**:
  - Periodically fetches new repositories from GitHub.
  - Sends repository fetched events to the Commits Manager Service.
  - Sends repository metadata events to the Commits Manager Service.
  - Retrieves the last repository fetch time via gRPC from the Commits Manager Service.
  
- **Responsibilities**:
  - Discover new repositories.
  - Ensure the Commits Manager Service is updated with the latest repository information.

### 3. Commits Monitor Service

- **Functionality**:
  - Periodically fetches new commits for all repositories from GitHub.
  - Sends commit fetched events to the Commits Manager Service.
  - Retrieves all repositories and the last commit fetch time via gRPC from the Commits Manager Service.
  
- **Responsibilities**:
  - Monitor and fetch new commits for existing repositories.
  - Ensure the Commits Manager Service is updated with the latest commit information.

## Implementation Details

### Commits Manager Service

- **REST APIs**:
  - For CRUD operations on repository and commit metadata.
  
- **Database**:
  - PostgreSQL tables for repositories and commits.
  
- **Message Broker**:
  - RabbitMQ to handle new repository and commit fetched events.
  
- **gRPC Server**:
  - To allow other microservices to query metadata.

### Repos Discovery Service

- **GitHub API Interaction**:
  - Fetches new repositories from GitHub.
  
- **Event Publishing**:
  - Publishes repository fetched events to RabbitMQ.
  
- **gRPC Client**:
  - Retrieves the last repository fetch time from the Commits Manager Service.

### Commits Monitor Service

- **GitHub API Interaction**:
  - Fetches new commits for repositories from GitHub.
  
- **Event Publishing**:
  - Publishes commit fetched events to RabbitMQ.
  
- **gRPC Client**:
  - Retrieves all repositories and the last commit fetch time from the Commits Manager Service.

### Data Storage

- **Repositories Table**:
  - Stores repository details such as name, description, URL, language, forks count, stars count, open issues count, watchers count, and creation/update dates.

- **Commits Table**:
  - Stores commit details such as SHA, URL, message, author name, author date, creation/update dates, and the associated repository name.

### Scheduling

- **Periodic Fetching**:
  - The Repos Discovery Service and Commits Monitor Service use Go's `time.Ticker` to schedule data fetching at regular intervals.
  
- **Duplicate Prevention**:
  - Ensures no duplicate commits by comparing fetched data with existing records in the database.

### Endpoints

- **List Repositories:**
    GET <http://localhost:8081/repositories>
    Retrieves a list of all repositories.

    Example

    ```bash
    curl http://localhost:8081/repositories?page=1&limit =10
    ```

- **Fetch Repository Commits:**
    GET <http://localhost:8081/commits/{repoName}>
    Retrieves commits for a specific repository.

    Example

    ```bash
    curl http://localhost:8081/commits/chromium?page=1&limit=10&startDate=2024-08-01T12:41:52Z&endDate=2024-08-01T12:52:26Z
    ```

- **Fetch Overall Top N Committers:**
    GET <http://localhost:8081/top-commit-authors?limit=10>
    Retrieves the top N commit authors overall.
- **Fetch Top N Committers for a Specific Repository:**
    GET <http://localhost:8081/top-commit-authors/{repoName}?limit=10>  
    Retrieves the top N commit authors for a specific repository.

    Example

    ```bash
    curl http://localhost:8081/top-commit-authors/chromium/?limit =10
    ```

### Unit Tests

- Unit tests are included to validate core functionalities,of data persistence in
commits-manager-service.
- To run unit tests
  - change you directory to commits-manager-service and run

    ```bash
    go test  -v -cover -short  ./...
    ```

## Setup and Usage

1. **Update Environment Variables:**

    - change you directory to project folder

    - Update `app.env` with your GitHub token and username:
    - Specify the START_DATE and END_DATE to fetch commits

    ```markdown
    DSN=host=postgres port=5432 user=postgres password=password dbname=github_tracker sslmode=disable timezone=UTC connect_timeout=5
    GITHUB_TOKEN=
    GITHUB_USERNAME=chromium
    START_DATE=2024-08-02T18:21:46Z
    END_DATE=2024-10-03T10:01:20Z
    ```

2. **Build and Run:**

    - Use Docker to build and start the services:

    ```markdown
    docker-compose up --build
    ```

Ensure to review the code, update the environment variables, and follow the provided instructions for a successful setup and execution of the service.
