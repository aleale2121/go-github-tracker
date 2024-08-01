CREATE TABLE repositories
(
    id BIGSERIAL PRIMARY KEY,
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
    id BIGSERIAL PRIMARY KEY,
    sha VARCHAR(255) UNIQUE,
    url VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    author_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    repository_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (repository_name) REFERENCES repositories(name) ON DELETE CASCADE
);

CREATE TABLE fetch_repos_metadata
(
    id BIGSERIAL PRIMARY KEY,
    total INT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE fetch_commits_metadata
(
    id BIGSERIAL PRIMARY KEY,
    repository_name VARCHAR(255) NOT NULL,
    total INT NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (repository_name) REFERENCES repositories(name)
);

