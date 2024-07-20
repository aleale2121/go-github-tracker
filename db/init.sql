CREATE TABLE repositories (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    url VARCHAR(255),
    language VARCHAR(255),
    forks_count INT,
    stars_count INT,
    open_issues_count INT,
    watchers_count INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE commits (
    sha VARCHAR(255) PRIMARY KEY,
    url VARCHAR(255),
    message TEXT,
    author_name VARCHAR(255),
    author_date TIMESTAMP
);
