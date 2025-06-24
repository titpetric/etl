# Test setup with git tags

The test `./git-tags.sh` implements a minimal sqlite driven test that
creates a database schema, inserts data, queries and updates it.

When you run the script, it will:

- Create a github_tags table,
- Print tables available,
- Insert some test data,
- Query what's imported,
- Update a record,
- Query after update.
