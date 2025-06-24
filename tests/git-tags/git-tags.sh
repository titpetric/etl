#!/bin/bash
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:git-tags.db"

# Clear initial db from any previous runs
rm -f git-tags.db

# Create schema
etl query git-tags.sql

# Print tables
etl tables | jq .

# Insert data into table
cat git-tags.json | etl insert github_tags

# Latest tags
etl query tags-latest.sql | jq .

# Tags with no created_at
etl query tags-no-date.sql | jq .

# Updating a known record
etl query tags-update.sql commit_sha=3f2e1d0c9b8a7e6d5f4c3b2a created_at=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Latest tags
etl query tags-latest.sql | jq .
