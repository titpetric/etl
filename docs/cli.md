# ETL CLI Reference

A quick reference for all ETL command-line operations.

## Installation

```bash
go install github.com/titpetric/etl/cmd/etl@main
```

Verify installation:
```bash
etl --help
```

## Configuration

Set environment variables to connect to your database:

### SQLite (Default)
```bash
export ETL_DB_DSN="sqlite://file:database.db"
```

### PostgreSQL
```bash
export ETL_DB_DSN="postgres://user:password@localhost:5432/dbname"
```

### MySQL
```bash
export ETL_DB_DSN="mysql://user:password@tcp(localhost:3306)/dbname"
```

## Creating Initial State

### Create a clean SQLite database with schema

```bash
# Set up the database
export ETL_DB_DSN="sqlite://file:myapp.db"

# Initialize schema from SQL file
etl query schema.sql
```

### Example schema.sql
```sql
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  total_amount DECIMAL(10, 2),
  status TEXT DEFAULT 'pending',
  order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Reading Data

### Get a single record

```bash
# Get one user by email
etl get users email=alice@example.com

# Get one order by ID
etl get orders id=1
```

### Get all records from a table

```bash
# Get all users
etl get users --all

# Get all orders
etl get orders --all
```

### Custom queries

```bash
# Query from file with parameter substitution
etl query queries/users-by-name.sql name="Alice"

# Multiple parameters
etl query queries/orders-by-user.sql user_id=1 status=pending
```

#### Example queries/users-by-name.sql
```sql
SELECT id, name, email, created_at
FROM users
WHERE name LIKE '%' || :name || '%'
ORDER BY created_at DESC;
```

## Writing Data

### Insert records

```bash
# Insert a single record with arguments
etl insert users name="Bob Smith" email="bob@example.com"

# Insert from JSON
echo '{"name":"Carol White","email":"carol@example.com"}' | etl insert users

# Insert from file with overrides
cat user.json | etl insert users status=active
```

### Update records

```bash
# Update by primary key
echo '{"name":"Alice Updated"}' | etl update users id=1

# Update with conditions
etl update users status=active created_at="2025-01-01"

# Bulk update from JSON
cat updates.json | etl update users id=:id
```

### Delete records (custom query approach)

Since the CLI doesn't provide direct delete functionality, create a query file:

#### queries/delete-user.sql
```sql
DELETE FROM users WHERE id = :id
RETURNING id, name, email;
```

Then execute:
```bash
etl query queries/delete-user.sql id=3
```

## Working with JSON Data

### Reading file data into columns

Use the `@file` syntax to read JSON from files:

```bash
# Store file contents as JSON column
etl insert documents file_data=@data.json
```

### Piping data through jq

Combine with `jq` for data transformation:

```bash
# Transform and insert
cat api-response.json | jq '.users[]' | etl insert users

# Filter and update
cat records.json | jq 'select(.status=="pending")' | etl update orders id=:id
```

## Server Mode

Start the API/Web server:

```bash
etl server
```

By default listens on the address configured in `etl.yml` (typically `:8080`).

See [docs/server.md](./server.md) for server configuration and API development.

## Examples

### Complete workflow: Create, read, update

```bash
# 1. Initialize database
export ETL_DB_DSN="sqlite://file:app.db"
etl query schema.sql

# 2. Insert users
echo '{"name":"Alice Johnson","email":"alice@example.com"}' | etl insert users
echo '{"name":"Bob Smith","email":"bob@example.com"}' | etl insert users

# 3. Get all users
etl get users --all

# 4. Update a user
echo '{"email":"alice.updated@example.com"}' | etl update users id=1

# 5. Get updated user
etl get users id=1

# 6. Start server for API access
etl server
```

## Tips & Tricks

- **Bulk operations**: Use pipes to process multiple records at once
- **Parameter formatting**: Values support dynamic data via shell expansion: `etl insert users created_at=$(date -I)`
- **Query files**: Keep `.sql` files database-specific to handle syntax differences
- **JSON validation**: Pipe through `jq` first to validate JSON before database operations
- **Dry runs**: Use `--dry` flag if available to preview operations without committing

## Troubleshooting

- **Connection errors**: Check `ETL_DB_DSN` environment variable
- **Missing table**: Run schema initialization with `etl query schema.sql`
- **JSON parse errors**: Validate JSON with `jq` before piping to etl
- **SQL syntax errors**: Test queries directly in database client first, then move to .sql files
