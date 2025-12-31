# Configuration

Config represents the overall configuration structure, which includes a server and multiple endpoints.

**Field: `Server` ([Server](#server))**
Server is used to configure the service.

**Field: `Storage` ([Storage](#storage))**
Storage configures the storage default.

**Field: `Include` (`[]string`)**
Include specifies additional config files to include and merge.

**Field: `Endpoints` ([[]*Endpoint](#endpoint))**
Endpoints contain a list of configured endpoints.

# Server

Server represents the server configuration with an address.

**Field: `HttpAddr` (`string`)**
HttpAddr contains the address the server should listen on. Example: ":3123".

**Field: `GrpcAddr` (`string`)**
GrpcAddr contains the address the server should listen on. Example: ":50051".

**Field: `Features` (`map[string]bool`)**
Features contains feature flags available for conditional query execution.

# Storage

Storage type configures database connection DSN.
The driver is automatically derived from the DSN connection string.

**Field: `DSN` (`string`)**
DSN configures the connection string for the database.
Supports mysql://, postgres://, postgresql://, sqlite://, and driver-specific formats.

# Endpoint

Endpoint represents an endpoint configuration with a path and handler.

**Field: `Path` (`string`)**
Path contains the path passed to the handler. It's a request path e.g. `/api`.

**Field: `Name` (`string`)**
Name describes what the endpoint returns.

**Field: `Methods` (`[]string`)**
Methods contains the HTTP request methods. If omitted, all methods are considered.

**Field: `Handler` ([Handler](#handler))**
Handler contains configuration related to the particular handler for the request.

# Handler

Handler represents the handler configuration for an endpoint.

**Field: `Type` (`string`)**
Type is the handler type to invoke for the endpoint.

**Field: `Command` (`[]string`)**
Command mimics os.Args for the handler.

**Field: `Library` (`string`)**
Library is an optional library parameter for a handler.

**Field: `Function` (`string`)**
Function is an optional function name parameter for a handler.

**Field: `Request` ([[]*Request](#request))**
Request specifies the upstream request handler path to invoke (for request type handlers).
Supports path parameters with brace syntax: /api/users/{id}

**Field: `Query` (`string`)**
Query is a single SQL query (deprecated, use Queries instead).

**Field: `Queries` ([[]*QueryDef](#querydef))**
Queries is a list of queries to execute in sequence (pipeline execution).

**Field: `Single` (`boolean`)**
Single returns a single object instead of array (applies to all queries).

**Field: `Parameters` (`map[string]`)**
Parameters are static parameters merged with request parameters.

**Field: `Transaction` ([Transaction](#transaction))**
Transaction configures transactional behavior for write operations.

**Field: `Cache` ([Cache](#cache))**
Cache configures response caching behavior.

**Field: `RateLimit` ([RateLimit](#ratelimit))**
RateLimit configures rate limiting for the endpoint.

**Field: `Response` ([Response](#response))**
Response configures the response format and headers.

**Field: `decoder` (`func`)**


# Request

Request configures the request format and headers.

**Field: `Path` (`string`)**
Path is the endpoint url to invoke.

**Field: `As` (`string`)**
As is the key for the decoded json value

**Field: `Method` (`string`)**
Method sets the request method. Default to GET.

**Field: `Headers` (`map[string]string`)**
Headers are custom HTTP request headers.

**Field: `Body` (`string`)**
Body is a custom HTTP request body.

# QueryDef

QueryDef represents a single query in a query pipeline.

**Field: `Query` (`string`)**
Query is the SQL query to execute.

**Field: `As` (`string`)**
As is the path where the result should be stored in the response.
If empty, result is merged into the current scope.

**Field: `If` (`string`)**
If is a condition expression evaluated against the current scope.
Query only executes if condition is true.

**Field: `For` (`string`)**
For is a loop expression to execute the query for each item.
Format: (idx, item) in items
Results are placed at the path specified in As.

# Transaction

Transaction configures transactional behavior for write operations.

**Field: `Enabled` (`boolean`)**
Enabled indicates whether transactions should be used.

**Field: `Retries` (`int`)**
Retries specifies the number of times to retry a failed transaction commit.

**Field: `RetryDelayMs` (`int`)**
RetryDelayMs specifies the delay in milliseconds between retry attempts.

# Cache

Cache configures response caching behavior.

**Field: `Enabled` (`boolean`)**
Enabled indicates whether response caching is enabled.

**Field: `Expire` (`string`)**
Expire specifies the cache time-to-live duration (e.g., "30s", "5m", "1h").
Supports all time.Duration formats. Defaults to 5 minutes if not specified.

**Field: `KeyPattern` (`string`)**
KeyPattern specifies the cache key pattern (supports path and query parameters).

# RateLimit

RateLimit configures rate limiting for the endpoint.
Example: rate: 100, per: "1m" means 100 requests per minute (default).
Per defaults to "1m" if not specified, so rate defines requests per minute.

**Field: `enabled` (`boolean`)**
Enabled indicates whether rate limiting is enabled.

**Field: `rate` (`int`)**
Rate specifies the number of requests allowed in the time period.
Defaults to requests per minute (when Per is unset).
Example: 100 = 100 requests per minute

**Field: `per` (`string`)**
Per specifies the time interval for the rate limit (e.g., "30s", "5m", "1h").
Supports all time.Duration formats. Defaults to "1m" (one minute) if not specified.
Example: "30s" means 30 seconds, "5m" means 5 minutes

# Response

Response configures the response format and headers.

**Field: `Headers` (`map[string]string`)**
Headers are custom HTTP response headers.
If not specified, Content-Type defaults to application/json for JSON responses
or text/html; charset=utf-8 for template responses.

**Field: `Template` (`string`)**
Template is a VueGo template string for formatting the response.
If specified, the response will be rendered using this template.

# Path

Path is the combination of Method and Path.

**Field: `Methods` (`[]string`)**
Methods contains the HTTP request methods. If omitted, all methods are considered.

**Field: `Path` (`string`)**
Path contains the request path for the router, e.g. `/users/{id}`.

