CREATE TABLE IF NOT EXISTS github_tags (
    name TEXT PRIMARY KEY,      -- Unique identifier, not null
    commit_sha TEXT NOT NULL,   -- Required field
    created_at TEXT             -- Nullable field
);
