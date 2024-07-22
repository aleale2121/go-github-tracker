CREATE TABLE repositories (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
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

CREATE TABLE commits (
    sha VARCHAR(255) PRIMARY KEY,
    url VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    author_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    repository_id BIGINT NOT NULL,
    FOREIGN KEY (repository_name) REFERENCES repositories(name)
);


CREATE TABLE fetch_metadata (
    id SERIAL PRIMARY KEY,
    last_repository_fetch TIMESTAMP,
    last_commit_fetch TIMESTAMP
);
