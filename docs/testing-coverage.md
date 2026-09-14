# Testing coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to
cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                             | Coverage | Cognitive | Lines |
|--------|-------------------------------------|----------|-----------|-------|
| ❌    | cmd/etl                             | 0.00%    | 7         | 69    |
| ❌    | drivers                             | 0.00%    | 35        | 237   |
| ❌    | handlers                            | 0.00%    | 106       | 431   |
| ❌    | internal                            | 0.00%    | 25        | 100   |
| ❌    | model                               | 0.00%    | 23        | 83    |
| ✅    | server                              | 20.69%   | 2         | 40    |
| ❌    | server/config                       | 62.50%   | 46        | 88    |
| ✅    | server/config/loader                | 90.64%   | 24        | 263   |
| ✅    | server/internal                     | 100.00%  | 2         | 14    |
| ✅    | server/internal/db/order            | 100.00%  | 0         | 12    |
| ✅    | server/internal/handler             | 80.77%   | 12        | 35    |
| ✅    | server/internal/handler/model       | 100.00%  | 2         | 15    |
| ❌    | server/internal/handler/query       | 3.95%    | 32        | 120   |
| ✅    | server/internal/handler/query/model | 71.43%   | 2         | 12    |
| ❌    | server/internal/handler/request     | 76.27%   | 36        | 173   |
| ❌    | server/internal/handler/sql         | 54.63%   | 180       | 698   |
| ✅    | server/middleware/cache             | 88.84%   | 60        | 288   |
| ✅    | server/middleware/ratelimit         | 84.54%   | 50        | 301   |

## Functions

| Status | Package                             | Function                         | Coverage | Cognitive |
|--------|-------------------------------------|----------------------------------|----------|-----------|
| ❌    | cmd/etl                             | HandleCommand                    | 0.00%    | 1         |
| ❌    |                                     | getInput                         | 0.00%    | 2         |
| ❌    |                                     | main                             | 0.00%    | 1         |
| ❌    |                                     | start                            | 0.00%    | 3         |
| ❌    | drivers                             | MySQL.Insert                     | 0.00%    | 6         |
| ❌    |                                     | MySQL.Query                      | 0.00%    | 2         |
| ❌    |                                     | MySQL.Tables                     | 0.00%    | 3         |
| ✅    |                                     | MySQL.insertQueryNamed           | 0.00%    | 0         |
| ❌    |                                     | New                              | 0.00%    | 1         |
| ✅    |                                     | NewMySQL                         | 0.00%    | 0         |
| ✅    |                                     | NewPgx                           | 0.00%    | 0         |
| ✅    |                                     | NewSqlite                        | 0.00%    | 0         |
| ❌    |                                     | Pgx.Insert                       | 0.00%    | 6         |
| ❌    |                                     | Pgx.Query                        | 0.00%    | 2         |
| ✅    |                                     | Pgx.Tables                       | 0.00%    | 0         |
| ✅    |                                     | Pgx.insertQueryNamed             | 0.00%    | 0         |
| ❌    |                                     | Sqlite.Insert                    | 0.00%    | 8         |
| ❌    |                                     | Sqlite.Query                     | 0.00%    | 5         |
| ✅    |                                     | Sqlite.Tables                    | 0.00%    | 0         |
| ❌    |                                     | Sqlite.insertQueryNamed          | 0.00%    | 2         |
| ❌    | handlers                            | Get                              | 0.00%    | 19        |
| ❌    |                                     | Insert                           | 0.00%    | 2         |
| ❌    |                                     | List                             | 0.00%    | 11        |
| ❌    |                                     | Query                            | 0.00%    | 12        |
| ✅    |                                     | Server                           | 0.00%    | 0         |
| ❌    |                                     | Tables                           | 0.00%    | 2         |
| ❌    |                                     | Update                           | 0.00%    | 24        |
| ❌    |                                     | UpdateRequest                    | 0.00%    | 7         |
| ❌    |                                     | Version                          | 0.00%    | 4         |
| ❌    |                                     | buildInsertQuery                 | 0.00%    | 3         |
| ❌    |                                     | buildUpdateQuery                 | 0.00%    | 8         |
| ✅    |                                     | csv                              | 0.00%    | 0         |
| ❌    |                                     | decodeQueryParameters            | 0.00%    | 10        |
| ❌    |                                     | scanAllRecords                   | 0.00%    | 3         |
| ❌    |                                     | scanRecord                       | 0.00%    | 1         |
| ❌    | internal                            | DecodeQuery                      | 0.00%    | 10        |
| ❌    |                                     | DecodeRecords                    | 0.00%    | 7         |
| ❌    |                                     | Scan                             | 0.00%    | 1         |
| ❌    |                                     | ScanAll                          | 0.00%    | 4         |
| ❌    |                                     | Statements                       | 0.00%    | 3         |
| ✅    |                                     | builtins                         | 0.00%    | 0         |
| ❌    | model                               | Config.ParseFlags                | 0.00%    | 1         |
| ✅    |                                     | NewConfig                        | 0.00%    | 0         |
| ✅    |                                     | NewFlagSet                       | 0.00%    | 0         |
| ❌    |                                     | RecordInput.Record               | 0.00%    | 1         |
| ❌    |                                     | dbValue                          | 0.00%    | 2         |
| ❌    |                                     | filterKnownArgs                  | 0.00%    | 19        |
| ✅    | server                              | Module.Mount                     | 100.00%  | 0         |
| ✅    |                                     | Module.Name                      | 100.00%  | 0         |
| ✅    |                                     | Module.Start                     | 100.00%  | 0         |
| ✅    |                                     | Module.Stop                      | 100.00%  | 0         |
| ✅    |                                     | NewConfig                        | 100.00%  | 0         |
| ✅    |                                     | NewModule                        | 100.00%  | 0         |
| ❌    |                                     | Start                            | 0.00%    | 2         |
| ❌    | server/config                       | Config.Validate                  | 0.00%    | 2         |
| ❌    |                                     | Decode                           | 66.67%   | 41        |
| ✅    |                                     | Handler.Decode                   | 66.67%   | 1         |
| ✅    |                                     | Handler.UnmarshalYAML            | 100.00%  | 0         |
| ✅    |                                     | applyStorageEnvOverrides         | 75.00%   | 2         |
| ✅    | server/config/loader                | CacheCloneManager.Get            | 100.00%  | 1         |
| ✅    |                                     | CacheCloneManager.String         | 100.00%  | 0         |
| ✅    |                                     | CacheCloneManager.loadAndSet     | 66.67%   | 2         |
| ✅    |                                     | CacheCloneManager.set            | 90.91%   | 1         |
| ✅    |                                     | CacheExpiryManager.Get           | 100.00%  | 2         |
| ✅    |                                     | CacheExpiryManager.String        | 100.00%  | 0         |
| ✅    |                                     | CacheExpiryManager.loadAndSet    | 66.67%   | 2         |
| ✅    |                                     | CacheExpiryManager.set           | 100.00%  | 0         |
| ✅    |                                     | CacheForeverManager.Get          | 100.00%  | 1         |
| ✅    |                                     | CacheForeverManager.String       | 100.00%  | 0         |
| ✅    |                                     | CacheForeverManager.loadAndSet   | 66.67%   | 2         |
| ✅    |                                     | CacheForeverManager.set          | 100.00%  | 0         |
| ✅    |                                     | CacheModifiedManager.Get         | 86.67%   | 5         |
| ✅    |                                     | CacheModifiedManager.String      | 100.00%  | 0         |
| ✅    |                                     | CacheModifiedManager.loadAndSet  | 66.67%   | 2         |
| ✅    |                                     | CacheModifiedManager.set         | 100.00%  | 0         |
| ✅    |                                     | CacheNone.Get                    | 75.00%   | 1         |
| ✅    |                                     | CacheNone.String                 | 100.00%  | 0         |
| ✅    |                                     | CacheNone.set                    | 0.00%    | 0         |
| ✅    |                                     | CacheSharedManager.Get           | 100.00%  | 1         |
| ✅    |                                     | CacheSharedManager.String        | 100.00%  | 0         |
| ✅    |                                     | CacheSharedManager.loadAndSet    | 66.67%   | 2         |
| ✅    |                                     | CacheSharedManager.set           | 90.91%   | 1         |
| ✅    |                                     | Load                             | 100.00%  | 1         |
| ✅    |                                     | Loader.Load                      | 100.00%  | 0         |
| ✅    |                                     | New                              | 100.00%  | 0         |
| ✅    |                                     | NewCacheCloneManager             | 100.00%  | 0         |
| ✅    |                                     | NewCacheExpiryManager            | 100.00%  | 0         |
| ✅    |                                     | NewCacheForeverManager           | 100.00%  | 0         |
| ✅    |                                     | NewCacheModifiedManager          | 100.00%  | 0         |
| ✅    |                                     | NewCacheNone                     | 100.00%  | 0         |
| ✅    |                                     | NewCacheSharedManager            | 100.00%  | 0         |
| ✅    | server/internal                     | Marshal                          | 100.00%  | 1         |
| ✅    |                                     | MarshalIndent                    | 100.00%  | 1         |
| ✅    | server/internal/db/order            | Asc                              | 100.00%  | 0         |
| ✅    |                                     | Desc                             | 100.00%  | 0         |
| ✅    | server/internal/handler             | Mount                            | 80.77%   | 12        |
| ✅    | server/internal/handler/model       | DBValue                          | 100.00%  | 2         |
| ✅    |                                     | Handlers                         | 100.00%  | 0         |
| ✅    |                                     | Register                         | 100.00%  | 0         |
| ❌    | server/internal/handler/query       | Handler.Handler                  | 0.00%    | 2         |
| ❌    |                                     | Handler.ServeHTTP                | 0.00%    | 3         |
| ✅    |                                     | Handler.Type                     | 100.00%  | 0         |
| ❌    |                                     | Handler.eval                     | 0.00%    | 12        |
| ❌    |                                     | Handler.prepareQueryParams       | 0.00%    | 15        |
| ✅    |                                     | NewHandler                       | 100.00%  | 0         |
| ✅    |                                     | init                             | 100.00%  | 0         |
| ✅    | server/internal/handler/query/model | Load                             | 71.43%   | 2         |
| ✅    | server/internal/handler/request     | Handler.EvaluateRequest          | 89.13%   | 14        |
| ✅    |                                     | Handler.Handler                  | 86.96%   | 3         |
| ❌    |                                     | Handler.ServeHTTP                | 50.00%   | 12        |
| ✅    |                                     | Handler.Type                     | 100.00%  | 0         |
| ❌    |                                     | Handler.buildUpstreamPath        | 60.00%   | 6         |
| ✅    |                                     | Handler.renderTemplateResponse   | 80.00%   | 1         |
| ✅    |                                     | NewHandler                       | 100.00%  | 0         |
| ✅    |                                     | init                             | 100.00%  | 0         |
| ✅    | server/internal/handler/sql         | Handler.Handler                  | 88.24%   | 17        |
| ❌    |                                     | Handler.ServeHTTP                | 76.60%   | 20        |
| ✅    |                                     | Handler.Type                     | 100.00%  | 0         |
| ✅    |                                     | Handler.buildCacheKey            | 100.00%  | 2         |
| ✅    |                                     | Handler.collectParameters        | 88.89%   | 15        |
| ❌    |                                     | Handler.evaluateCondition        | 0.00%    | 3         |
| ❌    |                                     | Handler.executeLoop              | 0.00%    | 13        |
| ❌    |                                     | Handler.executePipeline          | 69.57%   | 32        |
| ✅    |                                     | Handler.executeQuery             | 100.00%  | 2         |
| ✅    |                                     | Handler.executeQueryDirect       | 90.48%   | 9         |
| ✅    |                                     | Handler.executeQueryTx           | 80.95%   | 9         |
| ❌    |                                     | Handler.executeWithTransaction   | 48.00%   | 15        |
| ✅    |                                     | Handler.getFromCache             | 100.00%  | 2         |
| ❌    |                                     | Handler.renderTemplateResponse   | 0.00%    | 8         |
| ❌    |                                     | Handler.setAtPath                | 0.00%    | 6         |
| ❌    |                                     | Handler.setAtPathWithIndex       | 0.00%    | 9         |
| ✅    |                                     | Handler.setInCache               | 100.00%  | 4         |
| ✅    |                                     | Handler.setRateLimitHeaders      | 66.67%   | 4         |
| ❌    |                                     | Handler.setResponseHeaders       | 30.00%   | 8         |
| ✅    |                                     | Handler.shouldUseTransaction     | 100.00%  | 1         |
| ✅    |                                     | NewHandler                       | 100.00%  | 0         |
| ✅    |                                     | init                             | 100.00%  | 0         |
| ✅    |                                     | newTemplateFS                    | 0.00%    | 0         |
| ❌    |                                     | templateFS.Open                  | 0.00%    | 1         |
| ✅    |                                     | templateFile.Close               | 0.00%    | 0         |
| ✅    |                                     | templateFile.Read                | 0.00%    | 0         |
| ✅    |                                     | templateFile.Stat                | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.IsDir           | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.ModTime         | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.Mode            | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.Name            | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.Size            | 0.00%    | 0         |
| ✅    |                                     | templateFileInfo.Sys             | 0.00%    | 0         |
| ✅    | server/middleware/cache             | CustomKeyBuilder.BuildKey        | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.BuildKey       | 100.00%  | 2         |
| ✅    |                                     | DefaultKeyBuilder.WithHeaders    | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.WithPattern    | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.WithQuery      | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.extractHeaders | 100.00%  | 3         |
| ✅    |                                     | DefaultKeyBuilder.extractQuery   | 100.00%  | 12        |
| ❌    |                                     | MemoryStore.CleanupExpired       | 0.00%    | 4         |
| ✅    |                                     | MemoryStore.Clear                | 90.00%   | 1         |
| ✅    |                                     | MemoryStore.Delete               | 90.00%   | 1         |
| ✅    |                                     | MemoryStore.Get                  | 92.86%   | 3         |
| ✅    |                                     | MemoryStore.Set                  | 100.00%  | 2         |
| ❌    |                                     | Middleware.CleanupExpired        | 0.00%    | 1         |
| ✅    |                                     | Middleware.WithEnabled           | 100.00%  | 0         |
| ✅    |                                     | Middleware.WithLogger            | 0.00%    | 0         |
| ✅    |                                     | Middleware.WithTTL               | 100.00%  | 0         |
| ✅    |                                     | Middleware.Wrap                  | 95.83%   | 31        |
| ✅    |                                     | NewCustomKeyBuilder              | 100.00%  | 0         |
| ✅    |                                     | NewDefaultKeyBuilder             | 100.00%  | 0         |
| ✅    |                                     | NewMemoryStore                   | 100.00%  | 0         |
| ✅    |                                     | NewMiddleware                    | 100.00%  | 0         |
| ✅    | server/middleware/ratelimit         | CustomKeyBuilder.BuildKey        | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.BuildKey       | 100.00%  | 2         |
| ✅    |                                     | DefaultKeyBuilder.WithHeaders    | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.WithPattern    | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.WithQuery      | 100.00%  | 0         |
| ✅    |                                     | DefaultKeyBuilder.extractHeaders | 100.00%  | 3         |
| ✅    |                                     | DefaultKeyBuilder.extractQuery   | 100.00%  | 12        |
| ✅    |                                     | DefaultKeyBuilder.getClientIP    | 100.00%  | 5         |
| ❌    |                                     | MemoryStore.CleanupExpired       | 0.00%    | 4         |
| ✅    |                                     | MemoryStore.Clear                | 100.00%  | 1         |
| ✅    |                                     | MemoryStore.Inc                  | 89.47%   | 3         |
| ✅    |                                     | MemoryStore.Rate                 | 92.86%   | 3         |
| ✅    |                                     | MemoryStore.Reset                | 100.00%  | 1         |
| ❌    |                                     | Middleware.CleanupExpired        | 0.00%    | 1         |
| ✅    |                                     | Middleware.WithEnabled           | 100.00%  | 0         |
| ✅    |                                     | Middleware.WithLogger            | 0.00%    | 0         |
| ✅    |                                     | Middleware.Wrap                  | 86.49%   | 12        |
| ✅    |                                     | NewCustomKeyBuilder              | 100.00%  | 0         |
| ✅    |                                     | NewDefaultKeyBuilder             | 100.00%  | 0         |
| ✅    |                                     | NewMemoryStore                   | 100.00%  | 0         |
| ✅    |                                     | NewMiddleware                    | 100.00%  | 0         |
| ✅    |                                     | max                              | 100.00%  | 1         |
| ❌    |                                     | responseWrapper.Flush            | 0.00%    | 1         |
| ❌    |                                     | responseWrapper.Hijack           | 0.00%    | 1         |

