# ETL tooling for your CI pipelines

The etl cli enables you to interface your database from shell scripts.

It supports:

- sqlite,
- postgres,
- mysql

Check out [tests/git-tags.sh](./tests/git-tags.sh) for an usage example.

## Installation

```bash
go install github.com/titpetric/etl/cmd/etl@main
```

If cloning the repository run `task` in the project root.

## Usage examples

**To insert records**:

- `etl insert <table> [column=value column2=value2 ...]`
- `cat records.json | etl insert <table> [column=value ...]`

The input json is optional. Column values can be overriden as arguments.
The passed column value supports reading in files with `json=@<file>`.
This will set the JSON data to the column named `json` in the database.

To update records, use `etl update`.

**To get records**:

- `etl get <table> [--all] [column=value ...]`

It will return one record, unless `--all` is provided. The optional
column arguments are used to filter data for a `WHERE` clause.

**To use custom queries**:

- `etl query <file.sql> [column=value ...]`

The query file supports named parameters which are filled from arguments
(`:column`). The query files are usually coupled to the database you're
working with due to differences between SQL syntax.

> In order to implement the correct deletion or truncation behaviour for
> any database, a query file should be created. The cli doesn't provide
> any `delete` or `truncate` functionality for several unstated reasons.

## Motivation

I didn't find any nice tooling that would let me create and update
records in a set of desired database types (mysql, pgsql, sqlite).

The `etl` cli is an attempt to provide a database agnostic interface
that allows one to either use a sqlite db temporarily or connects to
persistent storage.

One could use database particular clients (`pgsql...`, `mysql -e`,...).

The `etl` tool does something similar in order to:

- provide better support for JSON data sources as input
- provide simple insert/update support with arguments over json
- provide query capability allowing customization

It's intended that `etl` is used in combination with `jq` or `yq` to
process JSON data before storing it into the database.

If you need proper database migrations, take a look at
[go-bridget/mig](https://github.com/go-bridget/mig).

## Configuring

All you need to do to make the cli functional is to declare the
following environment variables:

```bash
export ETL_DB_DRIVER=sqlite
export ETL_DB_DSN="file:git-tags.db"
```

For MySQL you would do:

```bash
export ETL_DB_DRIVER=mysql
export ETL_DB_DSN="etl:etl@tcp(localhost:3306)/etl"
```

For Postgres something like:

```bash
export ETL_DB_DRIVER=postgres
export ETL_DB_DSN="postgres://username:password@localhost:5432/database_name"
```
