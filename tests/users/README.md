# Users Test Configuration

This directory contains example configurations for ETL endpoints demonstrating key features with CQRS pattern separation.

## Files

- **etl.yml** - Main configuration that includes other configs
- **etl.users.yml** - User query endpoints (4 endpoints: 2 JSON + 2 HTML)
- **etl.users_write.yml** - User command endpoints (3 endpoints)
- **etl.login.yml** - Login query endpoints (1 endpoint)
- **etl.login_write.yml** - Login command endpoints (2 endpoints)
- **etl.orders.yml** - Order query endpoints (3 endpoints)
- **etl.orders_write.yml** - Order command endpoints (1 endpoint)

## Features Demonstrated

### 1. Include Directive
The main `etl.yml` uses the `include:` directive to load additional configuration files:
```yaml
include:
  - etl.users.yml
  - etl.users_write.yml
  - etl.login.yml
  - etl.login_write.yml
  - etl.orders.yml
  - etl.orders_write.yml
```

All endpoints from included files are merged into the main configuration. Endpoints are appended in include order.

### 2. Response Templates with VueGo
The `etl.users.yml` demonstrates HTML template responses alongside JSON:

```yaml
response:
  headers:
    "Content-Type": "text/html; charset=utf-8"
  template: |
    <h1>Users List</h1>
    <table border="1" cellpadding="8">
      ...
      <tr v-for="user in users">
        <td>{{ user.id }}</td>
        ...
      </tr>
      ...
    </table>
```

- **headers**: Custom HTTP response headers (optional)
- **template**: VueGo template string for rendering responses with Vue-like syntax

### 3. Rate Limiting Headers
All endpoints with `rateLimit` enabled automatically include:
- `X-RateLimit-Limit`: Maximum requests per second
- `X-RateLimit-Remaining`: Burst size

### 4. Caching
Endpoints demonstrate caching with duration-based invalidation:
```yaml
cache:
  enabled: true
  expire: "5m"
  keyPattern: "users:all"
```

## Running the Tests

Validate the configuration loads correctly:
```bash
cd tests/users
go test ./...
```

This runs the `TestValidateConfig` test which verifies:
- Configuration loads successfully with all includes
- Expected 14 endpoints are present (8 query + 6 command with CQRS pattern)
- All endpoints have valid paths and HTTP methods

## Endpoints Overview

### User Query Endpoints (etl.users.yml)
- `GET /api/users` - List all users (JSON)
- `GET /api/users/{id}` - Get user by ID (JSON)
- `GET /users` - List all users (HTML table)
- `GET /users/{id}` - Get user details (HTML card)

### User Command Endpoints (etl.users_write.yml)
- `POST /api/users` - Create new user
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user

### Order Query Endpoints (etl.orders.yml)
- `GET /orders` - List all orders (JSON)
- `GET /orders/{id}` - Get order by ID with user details (JSON)
- `GET /users/{user_id}/orders` - Get user's orders (JSON)

### Order Command Endpoints (etl.orders_write.yml)
- `POST /orders` - Create new order

### Login Query Endpoints (etl.login.yml)
- `POST /login` - User login

### Login Command Endpoints (etl.login_write.yml)
- `POST /logout` - User logout
- `POST /refresh-token` - Refresh auth token

## Configuration Merging

When multiple config files are included:

1. **Endpoints**: Appended in include order
2. **Server settings**: Later includes override earlier ones
3. **Storage config**: Later includes override earlier ones
4. **Features**: Later includes merge/override with earlier ones
