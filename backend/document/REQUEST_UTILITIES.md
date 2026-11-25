# Shared request utilities

The starter keeps cross-cutting code in small local packages.

| Concern | Implementation |
|---|---|
| Structured logs | `pkg/logger` configures Zap and adds the service name to each entry. |
| Request identifiers | `internal/middleware` creates a request ID and carries a caller's correlation ID forward. Both IDs are available in the Gin context and response headers. |
| Panic recovery | Gin recovery is installed during router setup in `cmd/main.go`. |
| Responses and errors | Handlers map use-case results and domain errors to HTTP status codes and JSON. Health checks have their own minimal response contract. |
