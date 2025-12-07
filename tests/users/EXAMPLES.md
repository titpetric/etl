# Response Template Examples

## JSON Response (Default)

When no response template is specified, responses are JSON encoded:

```yaml
- path: "/users"
  handler:
    type: "sql"
    query: "SELECT id, name, email FROM users"
```

Response:
```json
[
  {
    "id": 1,
    "name": "Alice Johnson",
    "email": "alice@example.com"
  },
  {
    "id": 2,
    "name": "Bob Smith",
    "email": "bob@example.com"
  }
]
```

## HTML Template Response

Use the `response:` section with a `template:` string to render HTML:

```yaml
- path: "/users-html"
  handler:
    type: "sql"
    query: "SELECT id, name, email FROM users"
    response:
      headers:
        "X-Custom-Header": "users-list"
      template: |
        <!DOCTYPE html>
        <html>
        <head><title>Users</title></head>
        <body>
          <h1>Users List</h1>
          <table>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Email</th>
            </tr>
            {{ range .data }}
            <tr>
              <td>{{ .id }}</td>
              <td>{{ .name }}</td>
              <td>{{ .email }}</td>
            </tr>
            {{ end }}
          </table>
        </body>
        </html>
```

Response: HTML content with table

## Template Data Context

The template receives data in the following ways:

### 1. Full Result as `.data`
```vue
{{ range .data }}
  {{ .field_name }}
{{ end }}
```

### 2. Individual Fields (for single results)
When `single: true`:
```vue
<h1>{{ .name }}</h1>
<p>Email: {{ .email }}</p>
```

### 3. Arrays in Template
```vue
{{ range .data }}
  <li>{{ .name }}</li>
{{ end }}
```

## Response Headers

Custom headers can be specified:

```yaml
response:
  headers:
    "Content-Type": "text/html; charset=utf-8"
    "X-Custom-Header": "custom-value"
    "Cache-Control": "no-cache"
  template: "..."
```

**Default Headers:**
- JSON response: `Content-Type: application/json`
- HTML response: `Content-Type: text/html; charset=utf-8`

If you specify a header with the same name, your value overrides the default.

## Rate Limit Headers

When rate limiting is enabled:

```yaml
rateLimit:
  enabled: true
  requestsPerSecond: 100
  burstSize: 200
```

Response automatically includes:
- `X-RateLimit-Limit: 100`
- `X-RateLimit-Remaining: 200`

## Example: User Detail Card

```yaml
- path: "/user/{id}/card"
  handler:
    type: "sql"
    single: true
    query: "SELECT id, name, email, created_at FROM users WHERE id = :id"
    response:
      template: |
        <!DOCTYPE html>
        <html>
        <head>
          <style>
            .card {
              border: 1px solid #ddd;
              border-radius: 4px;
              padding: 20px;
              max-width: 400px;
            }
          </style>
        </head>
        <body>
          <div class="card">
            <h2>{{ .name }}</h2>
            <p><strong>Email:</strong> <a href="mailto:{{ .email }}">{{ .email }}</a></p>
            <p><strong>Member Since:</strong> {{ .created_at }}</p>
          </div>
        </body>
        </html>
```

## VueGo Template Syntax

The template uses VueGo syntax which is similar to Go's text/template:

```vue
<!-- Variable substitution -->
{{ .fieldName }}

<!-- Range loops -->
{{ range .items }}
  <li>{{ .name }}</li>
{{ end }}

<!-- Conditionals -->
{{ if .is_active }}
  <span class="active">Active</span>
{{ end }}

<!-- Pipes and functions -->
{{ .created_at | formatDate }}
```

## Complete Configuration Example

See `tests/users/users-html.yml` for complete working examples with:
- HTML list endpoint
- HTML detail endpoint with styling
- Custom headers
- Rate limiting
- Caching
