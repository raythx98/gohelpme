# Modules

Per-package summaries — what each package owns and its public interface.

---

## `middleware/chain.go`

Middleware composition.

**Exports:**
- `Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc` — wraps a handler with middlewares in reverse order (last listed executes first/outermost).

---

## `middleware/log.go`

Structured request/response logging using zerolog.

**Exports:**
- `Log(logger ILogger, config LogConfig) Middleware` — logs method, path, status, latency, and request ID. Redacts configured fields.

---

## `middleware/errorhandler.go`

Maps application errors to HTTP responses.

**Exports:**
- `ErrorHandler(logger ILogger) Middleware` — converts `AppError`, `AuthError`, and validation errors to JSON responses with appropriate HTTP status codes.

---

## `middleware/jwtauth.go`

JWT Bearer token authentication.

**Exports:**
- `JwtAuth(jwt IJwt) Middleware` — validates Bearer token; rejects with 401 on failure.

---

## `middleware/basicauth.go`

HTTP Basic Auth authentication.

**Exports:**
- `BasicAuth(username, password string) Middleware` — validates Basic credentials; rejects with 401 on failure.

---

## `middleware/jwtorbasicauth.go`

Combined auth: tries Bearer first, falls back to Basic Auth.

**Exports:**
- `JwtOrBasicAuth(jwt IJwt, username, password string) Middleware`

---

## `middleware/ratelimit.go`

Token bucket rate limiting per IP or user.

**Exports:**
- `RateLimit(config RateLimitConfig) Middleware` — applies rate limits from config; returns 429 on excess.

---

## `middleware/cors.go`

CORS headers.

**Exports:**
- `Cors(allowedOrigins []string) Middleware`

---

## `middleware/recoverer.go`

Panic recovery. Logs the stack trace and returns 500.

**Exports:**
- `Recoverer(logger ILogger) Middleware`

---

## `middleware/redactor.go`

Sensitive field redaction in request/response logs.

**Exports:**
- `NewRedactor(fields []string) *Redactor` — creates a redactor for the given JSON field paths.
- `Redactor.Redact(body []byte) []byte` — replaces sensitive field values with `[REDACTED]`.

---

## `middleware/addrequestid.go`

Injects a UUID request ID into the request context and `X-Request-ID` response header.

**Exports:**
- `AddRequestId() Middleware`

---

## `middleware/reqctx.go`

Sets up the request context value container used by all other middleware.

**Exports:**
- `ReqCtx() Middleware` — must be the first middleware in the chain after CORS.

---

## `tool/logger/`

Structured logging interface.

**Exports:**
- `ILogger` interface: `Info(ctx, msg, fields...)`, `Warn(...)`, `Debug(...)`, `Error(...)`, `Fatal(...)`.
- `NewDefault(config LogConfig) ILogger` — creates a zerolog/slog JSON logger.
- `Field` type for structured log fields.

---

## `tool/jwthelper/`

JWT creation and validation.

**Exports:**
- `IJwt` interface: `CreateToken(claims Claims) (string, error)`, `ValidateToken(token string) (*Claims, error)`, `ExtractBearer(headers http.Header) string`.
- `New(secret string, ttl time.Duration) IJwt`

---

## `tool/crypto/`

Password hashing using Argon2id.

**Exports:**
- `ICrypto` interface: `Hash(password string) (string, error)`, `Verify(hash, password string) (bool, error)`.
- `New() ICrypto` — uses default Argon2id parameters (secure defaults).

---

## `tool/validator/`

Input validation wrapping `go-playground/validator`.

**Exports:**
- `IValidator` interface: `Validate(v interface{}) error`.
- `New() IValidator`

---

## `tool/postgres/`

PostgreSQL connection pool via pgx.

**Exports:**
- `IPostgres` interface: `Pool() *pgxpool.Pool`.
- `New(config IConfig) (IPostgres, error)` — creates a pgxpool connection.

---

## `tool/aws/`

AWS S3 integration.

**Exports:**
- `IS3` interface: `GetPresignedUrl(bucket, key string, expiry time.Duration) (string, error)`.
- `NewS3(cfg aws.Config, bucket string) IS3`

---

## `tool/reqctx/`

Request context value helpers.

**Exports:**
- `GetValue(ctx context.Context) *RequestContext` — returns the request-scoped value container.
- `RequestContext` struct: `RequestID string`, `UserID string`, `Error error` (and setters).

---

## `tool/random/`

Random string generation.

**Exports:**
- `AlphaNumeric(length int) string` — generates a cryptographically random alphanumeric string.

---

## `errorhelper/apperror.go`

Application error type.

**Exports:**
- `AppError` struct: `Code int`, `Message string`, `Err error`.
- `NewAppError(code int, message string, err error) *AppError`
- `IsAppError(err error) (*AppError, bool)`

---

## `errorhelper/autherror.go`

Authentication error type.

**Exports:**
- `AuthError` struct: `Message string`.
- `NewAuthError(message string) *AuthError`
- `IsAuthError(err error) (*AuthError, bool)`

---

## `errorhelper/dto.go`

JSON error response shapes returned to API clients.

---

## `builder/httprequest/` and `builder/httpclient/`

Fluent HTTP request builder.

**Exports:**
- `builder.New(ctx, method, url) *Builder` — creates a new request builder.
- `(*Builder).WithBody(data)`, `.WithHeaders(headers)`, `.WithTimeout(d)` — chainable configuration.
- `(*Builder).Build() (*http.Response, error)` — executes the request.
