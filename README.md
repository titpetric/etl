# ETL - SQL-First API & Web Development Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/titpetric/etl)](https://goreportcard.com/report/github.com/titpetric/etl)
[![Testing coverage](https://img.shields.io/badge/coverage-36.3%25-yellowgreen)](./docs/testing-coverage.md)
[![GoDoc](https://pkg.go.dev/badge/github.com/titpetric/etl.svg)](https://pkg.go.dev/github.com/titpetric/etl)

**ETL** is a no-code SQL-oriented framework.

ETL enables rapid development of APIs, web applications, and supports
CQRS design.

ETL allows you to define API endpoints declaratively in YAML, execute
parameterized SQL queries, and serve JSON APIs or render HTML
templates - all without writing Go code.

## Installation

```bash
go install github.com/titpetric/etl/cmd/etl@main
```

## Quick Navigation

- **[CLI Reference](./docs/cli.md)** – Command-line tool usage, database operations, and examples
- **[Configuration](./docs/config.md)** – Server configuration, endpoints, handlers, and options
- **[Server Development](./docs/server.md)** – API configuration, web development, and template rendering

## Overview

ETL supports these major use cases:

1. **CLI Database Interface** - Query and manipulate databases from shell scripts
2. **SQL-First API Development** - Expose database queries as json based REST/RPC API endpoints with zero boilerplate
3. **Web Development** - Responses can be templated, reusing existing APIs defined in the system as the data source
4. **CQRS Development** - Composable configuration allows trivial separation of read/write statements

## Supported Databases

- SQLite
- PostgreSQL
- MySQL

The ETL CLI aims to be database agnostic. It leans into `json` as a
portable data format, and implements compatibility drivers to talk to
each database, storing and returning JSON with a common API.

The SQL you define on your own is not portable, as it depends on the
database server in use. This allows you to use server specific syntax.

## Testing

Several test suites demonstrate core functionality:

- **[tests/git-tags](./tests/git-tags)** – CLI usage examples
- **[tests/petstore](./tests/petstore)** – SQL-as-API examples
- **[tests/users](./tests/users)** – CQRS and integrated testing

## Configuration

Set environment variables:

```bash
# SQLite (default)
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:myapp.db"

# PostgreSQL
export ETL_DB_DRIVER=postgres
export ETL_DB_DSN="user:pass@localhost:5432/dbname"

# MySQL
export ETL_DB_DRIVER=mysql
export ETL_DB_DSN="user:pass@tcp(localhost:3306)/dbname"
```

## Why ETL?

**Standard database clients** (`psql`, `mysql -e`) require mastering database-specific syntax and CLI arguments and flags.

**ETL provides:**

- Unified interface across SQLite, PostgreSQL, and MySQL
- First-class JSON support for input/output
- Simple insert/update operations via JSON or arguments
- Custom SQL query capability for complex operations
- Rapid API development without boilerplate
- Web application rendering alongside API servers

**Use with** `jq`, `yq`, or other standard tools for data transformation pipelines.

## License

[See LICENSE file](./LICENSE)

## References

- [go-bridget/mig](https://github.com/go-bridget/mig) – For database migrations
- [titpetric/vuego](https://github.com/titpetric/vuego) – Template engine
- [titpetric/platform](https://github.com/titpetric/platform) – HTTP server framework
