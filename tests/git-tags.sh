#!/bin/bash
export DB_DRIVER=sqlite
export DB_DSN="file:git-tags.db"

rm -f git-tags.db

# Create schema
etl query testdata/migrate.sql

# Print tables
etl tables | jq .

# Insert data into table
cat git-tags.json | etl insert github_tags

# Latest tags
etl query testdata/tags-latest.sql | jq .

# Tags with no created_at
etl query testdata/tags-no-date.sql | jq .

# Updating a known record
etl query testdata/tags-update.sql commit_sha=3f2e1d0c9b8a7e6d5f4c3b2a created_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Latest tags
etl query testdata/tags-latest.sql | jq .

# Truncate
etl truncate github_tags
