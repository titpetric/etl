# ETL - SQL-First API & Web Development Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/titpetric/etl)](https://goreportcard.com/report/github.com/titpetric/etl)
[![Coverage](https://img.shields.io/badge/coverage-31.7%25-orange)](./coverage.out)
[![GoDoc](https://pkg.go.dev/badge/github.com/titpetric/etl.svg)](https://pkg.go.dev/github.com/titpetric/etl)

**ETL** is a no-code SQL-oriented framework enabling rapid development of APIs, web applications, and CQRS systems. Define endpoints declaratively in YAML, execute parameterized SQL queries, and serve JSON APIs or render HTML templates—all without writing Go code.

## Quick Navigation

- **[CLI Reference](./docs/cli.md)** – Command-line tool usage, database operations, and examples
- **[Server Development](./docs/server.md)** – API configuration, web development, and template rendering

## Overview

ETL supports three major use cases:

1. **CLI Database Interface** – Query and manipulate databases from shell scripts
2. **SQL-First API Development** – Expose database queries as REST endpoints with zero boilerplate
3. **Web & CQRS Development** – Separate read and write models with composable HTTP-based architecture

### Supported Databases

- SQLite
- PostgreSQL
- MySQL

## Getting Started

### Installation

```bash
go install github.com/titpetric/etl/cmd/etl@main
```

### Quick Start with CLI

```bash
# Configure database
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:app.db"

# Initialize schema
etl query schema.sql

# Insert records
echo '{"name":"Alice","email":"alice@example.com"}' | etl insert users

# Query data
etl get users --all

# Update records
echo '{"name":"Alice Updated"}' | etl update users id=1
```

### Quick Start with Server

Create `etl.yml`:
```yaml
endpoints:
  - path: /api/users
    methods: [GET]
    handler:
      type: sql
      query: SELECT id, name, email FROM users ORDER BY id
```

Start server:
```bash
etl server
```

Access: `curl http://localhost:8080/api/users`

## Documentation

### [CLI Reference](./docs/cli.md)

Complete guide to command-line operations:
- Database configuration (SQLite, PostgreSQL, MySQL)
- Creating and initializing schemas
- Reading data (single records, all records, custom queries)
- Writing data (insert, update, delete with SQL files)
- Working with JSON data
- Full workflow examples

### [Server Development](./docs/server.md)

Build APIs and web applications with ETL:
- Simple API endpoints from SQL queries
- Web development with Vuego templates
- Composing APIs with external data sources
- CQRS patterns
- HTML rendering

## Testing

Two main test suites demonstrate core functionality:

- **[tests/git-tags](./tests/git-tags)** – CLI usage examples
- **[tests/petstore](./tests/petstore)** – SQL-as-API examples
- **[tests/users](./tests/users)** – CQRS and integrated testing

Run tests:
```bash
task up          # Start services
task test        # Run all tests
```

## Configuration

Set environment variables:

```bash
# SQLite (default)
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:myapp.db"

# PostgreSQL
export ETL_DB_DRIVER=postgres
export ETL_DB_DSN="postgres://user:pass@localhost:5432/dbname"

# MySQL
export ETL_DB_DRIVER=mysql
export ETL_DB_DSN="user:pass@tcp(localhost:3306)/dbname"
```

## Development

Requires [Task](https://taskfile.dev) for development workflows.

```bash
task setup    # Install development dependencies
task          # Format and build
task -l       # List all available tasks
```

## Why ETL?

**Standard database clients** (`psql`, `mysql -e`) require mastering database-specific syntax and CLIs.

**ETL provides:**
- Unified interface across SQLite, PostgreSQL, and MySQL
- First-class JSON support for input/output
- Simple insert/update operations via JSON or arguments
- Custom SQL query capability for complex operations
- Rapid API development without boilerplate
- Web application rendering alongside API servers

**Use with** `jq`, `yq`, or other standard tools for data transformation pipelines.

## Example: Complete Workflow

```bash
# 1. Create and initialize database
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:app.db"
etl query schema.sql

# 2. Insert bulk data from JSON file
cat users.json | jq '.[]' | etl insert users

# 3. Query and filter
etl query queries/active-users.sql | jq '.[] | select(.status=="active")'

# 4. Update records
cat updates.json | etl update users id=:id

# 5. Serve as API
etl server
```

## Architecture

```
┌─────────────┐
│   etl CLI   │  Database operations, custom queries
└──────┬──────┘
       │
   ┌───┴──────────────┐
   │   Databases      │
   ├──────────────────┤
   │ SQLite / PG / MY │
   └──────────────────┘
       
┌─────────────────────┐
│   etl server        │  API & Web
├─────────────────────┤
│  YAML Config        │
│  Routes → Handlers  │
│  SQL Queries        │
│  Templates (Vuego)  │
└────────┬────────────┘
         │
    ┌────┴─────────────────┐
    │  API Responses (JSON) │
    │  HTML Rendering      │
    └──────────────────────┘
```

## License

[See LICENSE file](./LICENSE)

## References

- [go-bridget/mig](https://github.com/go-bridget/mig) – For database migrations
- [titpetric/vuego](https://github.com/titpetric/vuego) – Template engine
- [titpetric/platform](https://github.com/titpetric/platform) – HTTP server framework
