# Architecture

## Overview

gohelpme is a **shared Go utility library** for personal microservices (go-dutch, go-zap, and others). It provides reusable middleware, HTTP client builders, authentication helpers, database connectivity, structured logging, and error handling patterns. It is a library — not a runnable service.

## Package Map

```
gohelpme/
├── middleware/           # net/http middleware implementations
│   ├── chain.go         # Middleware composition
│   ├── log.go           # Request/response structured logging
│   ├── errorhandler.go  # Error response mapping
│   ├── jwtauth.go       # JWT Bearer authentication
│   ├── basicauth.go     # HTTP Basic Auth
│   ├── jwtorbasicauth.go# Combined JWT-or-Basic-Auth
│   ├── ratelimit.go     # Token bucket rate limiting
│   ├── cors.go          # CORS headers
│   ├── recoverer.go     # Panic recovery
│   ├── redactor.go      # Sensitive field redaction in logs
│   ├── addrequestid.go  # Request ID injection
│   ├── reqctx.go        # Request-scoped context value management
│   └── example.go       # Usage examples
├── tool/                 # Utility packages
│   ├── logger/          # ILogger interface + zerolog/slog implementation
│   ├── jwthelper/       # IJwt interface + JWT creation/validation
│   ├── crypto/          # Password hashing (Argon2id)
│   ├── basicauth/       # Basic auth encoding/decoding
│   ├── validator/       # IValidator interface + go-playground wrapper
│   ├── postgres/        # IPostgres interface + pgx pool
│   ├── aws/             # AWS S3 integration
│   ├── httphelper/      # Generic HTTP client utilities
│   ├── reqctx/          # Request context value helpers
│   ├── random/          # Random string generation
│   ├── timehelper/      # Time utilities
│   └── inthelper/       # Integer utilities
├── builder/              # Fluent API builders
│   ├── httpclient/      # HTTP client builder
│   └── httprequest/     # HTTP request builder
├── errorhelper/          # Custom error types
│   ├── apperror.go      # AppError with code + message
│   ├── autherror.go     # AuthError for auth failures
│   ├── dto.go           # Error response DTOs
│   └── errorhelper.go   # Error utility functions
└── mocks/               # Mockery-generated mock implementations
```

## Design Principles

### Interface-First

Every major component defines an interface. Consumers depend on the interface, never the concrete type. This enables:
- Unit testing with mockery-generated mocks in consumer services
- Swapping implementations without changing consumer code
- Clear contract documentation

```go
// Consumer uses:
type Tools struct {
    Logger  logger.ILogger
    Crypto  crypto.ICrypto
    JWT     jwthelper.IJwt
}

// Not:
type Tools struct {
    Logger  *zerolog.Logger   // concrete — avoid
}
```

### Middleware Chain

Middleware follows a standard functional pattern:

```go
type Middleware func(next http.HandlerFunc) http.HandlerFunc

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    // applies in reverse order: last listed = outermost wrapper
}
```

### Request Context

Request-scoped values (request ID, error, user ID) are passed via the `reqctx` package rather than package-level globals or direct context keys:

```go
reqctx.GetValue(ctx).RequestID
reqctx.GetValue(ctx).SetError(err)
```

### Builder Pattern

HTTP clients use a fluent builder for constructing requests:

```go
builder.New(ctx, httprequest.Post, url).
    WithBody(data).
    WithHeaders(headers).
    Build()
```

## Usage in Consumer Services

Consumer services (go-dutch, go-zap) use this library for:
1. Middleware chain setup in their route registration
2. `Tools` struct initialization in their `resources/` package
3. `ILogger`, `ICrypto`, `IJwt` interfaces for dependency injection
4. Error types from `errorhelper/` for consistent error responses
5. Mockery-generated mocks for unit testing
