# Users Test Configuration

This directory contains example configurations for ETL endpoints demonstrating key features:

## Files

- **etl.yml** - Main configuration that includes other configs
- **users.yml** - User management endpoints (3 endpoints)
- **users-html.yml** - HTML template response examples (2 endpoints)
- **orders.yml** - Order management endpoints (4 endpoints)

## Features Demonstrated

### 1. Include Directive
The main `etl.yml` uses the `include:` directive to load additional configuration files:
```yaml
include:
  - "users.yml"
  - "users-html.yml"
  - "orders.yml"
```

All endpoints from included files are merged into the main configuration. Endpoints are appended in include order.

### 2. Response Templates with VueGo
The `users-html.yml` demonstrates HTML template responses:

```yaml
response:
  headers:
    "X-Custom-Header": "user-detail"
  template: |
    <!DOCTYPE html>
    <html>
    ...
    {{ range .data }}
      {{ .field_name }}
    {{ end }}
    ...
    </html>
```

- **headers**: Custom HTTP response headers (optional)
- **template**: VueGo template string for rendering responses

### 3. Rate Limiting Headers
All endpoints with `rateLimit` enabled automatically include:
- `X-RateLimit-Limit`: Maximum requests per second
- `X-RateLimit-Remaining`: Burst size

### 4. Caching
Endpoints demonstrate caching with TTL-based invalidation:
```yaml
cache:
  enabled: true
  ttlSeconds: 300
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
- Expected 14 endpoints are present (8 read + 6 write with CQRS pattern)
- All endpoints have valid paths and HTTP methods

## Endpoints Overview

### JSON Responses
- `GET /users` - List all users (JSON)
- `GET /users/{id}` - Get user by ID (JSON)
- `POST /users` - Create new user
- `GET /orders` - List all orders (JSON)
- `GET /orders/{id}` - Get order by ID with user details (JSON)
- `GET /users/{user_id}/orders` - Get user's orders (JSON)
- `POST /orders` - Create new order

### HTML Template Responses
- `GET /users-html` - List all users (HTML table)
- `GET /users-html/{id}` - Get user details (HTML card)

## Configuration Merging

When multiple config files are included:

1. **Endpoints**: Appended in include order
2. **Server settings**: Later includes override earlier ones
3. **Storage config**: Later includes override earlier ones
4. **Features**: Later includes merge/override with earlier ones
