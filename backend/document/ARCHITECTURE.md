# Architecture Blueprint

Derived from `backend-for-frontend`, `push-notification`, `sms-service`,
`email-service`, `email-gateway`, `as400-adapter`, `face-verification-service`,
`ekyc-service`. All 8 are Go + Gin HTTP services with the same skeleton:

```
cmd/            main.go — wiring only (config load, logger init, router, graceful shutdown)
internal/       everything private to this service
  handler/      Gin handlers — HTTP in/out only, no business logic
  <domain>/     one package per business capability (see "Domain packages" below)
  models/       request/response/domain structs
  middleware/   request IDs and other HTTP middleware
  routes/       route registration, grouped by prefix
  server/       http.Server setup, TLS, timeouts
config/         config.yaml + env-specific overrides (config.docker-compose.yaml, config.integration-test.yaml, ...)
pkg/            reusable service-owned code such as the Zap logger
bruno/          API docs + integration tests (see MOCKING_GUIDE.md)
document/       architecture notes, API docs, deployment handoff, release notes
db/             migrations + sqlc queries (only services with a relational DB)
test/           integration test fixtures (Mockoon datasources, scenarios)
vendor/         vendored dependencies used by the Docker build
```

## Request flow

### Clean layered architecture

**Routes → Handler → Actor → Aggregator → Client**

| Layer | Responsibility |
|---|---|
| Routes | Register endpoints and attach middleware. |
| Handler | Bind and validate HTTP input, invoke the actor, and map the result to an HTTP response. |
| Actor | Execute the use case and apply business rules. |
| Aggregator | Coordinate downstream client calls and combine their results for the actor. |
| Client | Handle downstream I/O and translate protocol-specific requests, responses, and errors. |

The flow shows the full call chain. A use case with a single downstream
operation can call its client directly when there is no coordination to do.
Wire the layers in `cmd/main.go` and inject the dependencies into their callers.

### HTTP request lifecycle

```
HTTP request
  → middleware (request-id, structured logging, panic recovery, auth)
  → routes/ dispatch to handler/
  → handler/ binds + validates request and calls the actor
  → actor applies business rules and calls the aggregator where needed
  → client performs downstream I/O (HTTP / DB / Kafka / SDK)
  → actor returns a result or a domain error
  → handler/ maps that to an HTTP response using the service's JSON contract
```

The handler layer is intentionally thin. Business logic, branching, and
cross-service calls belong in the domain package underneath it, not in the
`gin.HandlerFunc`. This is the strongest convention across all 8 services and
scales from BFF (18 domains, orchestrates many downstreams) down to
sms-service (one domain, one downstream).

## Domain packages (BFF's "actor" pattern)

BFF names its domain packages `internal/actor/<domain>` (e.g. `actor/payment`,
`actor/onboarding`) — each one owns calling the real downstream service(s) for
that business capability, with a matching `handler/<domain>` and
`models/<domain>`. Smaller services (sms-service, email-service) collapse this
into fewer, flatter packages since they only have one real capability. The
principle is the same at any scale: **one business capability = one package
that owns its downstream calls**, regardless of how many packages that turns
out to be.

## The adapter pattern (mandatory for anything external)

Every SDK-based external dependency — AWS (S3, STS), Kafka, Firebase, the
AS400 TCP link — must sit behind an interface defined in `internal/`, with the
real SDK call in one implementation and a stub/fake implementation for tests.
This is not optional: it's the only way `MOCKING_GUIDE.md`'s stub-mode testing
works for anything that isn't a plain HTTP call (HTTP downstreams get
Mockoon instead, which doesn't require an adapter in the same way).

```go
// internal/storage/storage.go
type Storage interface {
    Put(ctx context.Context, key string, body io.Reader) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
}

// internal/storage/s3.go        — real implementation, wraps aws-sdk-go-v2/service/s3
// internal/storage/stub.go      — in-memory fake, used when STUB_MODE=true or in unit tests
```

Construct the real-vs-stub choice once, in `cmd/main.go` or a `wire.go`, based
on config/env — never branch on "am I in a test" inside business logic.

## Health checks

Every service exposes the same three endpoints (see `bruno/Health/`):

| Path | Purpose |
|---|---|
| `GET /health/liveness` | process is up — used by Docker `HEALTHCHECK` and k8s liveness probe |
| `GET /health/readiness` | dependencies (DB, redis, kafka) are reachable — used by k8s readiness probe |
| `GET /ping` | trivial reachability check, no dependency checks |

## Configuration

`spf13/viper` loads a base `config.yaml` plus an environment-specific
override, selected via the `CONFIG_FILE` env var (e.g.
`CONFIG_FILE=/app/config/config.docker-compose.yaml`). Secrets are never
committed — they're injected as env vars or mounted files at deploy time, per
`CICD_GUIDE.md`.
