# Tech Stack

Versions shown are what's actually in the 8 services' `go.mod` today
(late 2025 generation). Pin to the newest of these when starting a new
service, not to any single service's exact version.

## Language & runtime

| | |
|---|---|
| Language | Go 1.26 |
| Build | `CGO_ENABLED=0`, static binary, `-ldflags="-w -s"` |
| Container base | `public.ecr.aws/docker/library/golang:1.26-alpine` (build), `public.ecr.aws/docker/library/alpine:3.21` (runtime) |
| Dependency mode | **vendored** (`-mod=vendor`) — mandatory, see "Why vendoring" below |

## Web framework

- **`gin-gonic/gin`** (v1.11–1.12 across services) — HTTP router + middleware chain. Every service uses it; not up for debate.

## Config, logging, validation

| Concern | Library | Notes |
|---|---|---|
| Config loading | `spf13/viper` | base `config.yaml` + env override via `CONFIG_FILE` |
| Structured logging | `go.uber.org/zap` via `pkg/logger` | don't call zap directly in handlers — use the local logger so log shape is consistent across services |
| Request validation | `go-playground/validator/v10` | struct tag validation on request-bound structs |
| Error types | Go errors | map domain errors to HTTP responses in handlers |
| Response envelope | Service-defined JSON | for example, `status` / `service_id` / `message` / `data` |
| Context propagation | Request context | request-id, trace-id helpers |
| Common middleware | Local Gin middleware | request IDs and Gin recovery |

## Data & infra clients

| Concern | Library | Used by |
|---|---|---|
| Relational DB | `jackc/pgx/v5` + `sqlc` + `golang-migrate`-style `db/migrations` | email-service, face-verification-service, ekyc-service |
| Cache / session / distributed lock | `redis/go-redis/v9` (or Valkey, wire-compatible) | BFF, face-verification-service, ekyc-service |
| Event streaming | `segmentio/kafka-go` (+ `sasl/aws_msk_iam_v2` for MSK IAM auth) | BFF (producer-only) |
| Object storage | `aws-sdk-go-v2/service/s3` (MinIO-compatible locally) | BFF |
| Push notification | `firebase.google.com/go/v4` | BFF, push-notification |
| Date/time | `dromara/carbon/v2` | BFF |
| Serialization | `vmihailenco/msgpack/v5` (messagepack) alongside standard `encoding/json` | BFF |
| Property-based testing | `leanovate/gopter` | BFF |
| Non-HTTP downstream (AS400) | raw TCP socket, custom framed message protocol (`internal/models/message`) | as400-adapter only |

## Tooling

| Tool | Purpose |
|---|---|
| `golangci-lint` v2 | linting — see `.golangci.yml` in `template-service/`, migrated from BFF's v1 ruleset (ekyc-service already did this migration; v1 config format is deprecated upstream) |
| `make` | single entrypoint for every dev action — `make help` lists targets |
| Bruno (`@usebruno/cli`) | API docs + integration tests, see MOCKING_GUIDE.md |
| Mockoon | HTTP downstream mocking for integration tests, opt-in — see MOCKING_GUIDE.md |
| Docker Compose | local dev stack (app + whichever infra it needs) |
| GitLab CI | build, lint, test, docker build/push, deploy — see CICD_GUIDE.md |

## Why vendoring

Run `make tidy vendor` in a generated service before building its Docker
image. The Docker build reads from `vendor/` with `-mod=vendor`, so Go
dependencies are resolved before the image build starts.

## Module naming

Pass the intended repository path to `--module`, for example
`github.com/your-org/my-service`. For a local example, use
`example.com/my-service`. The generator updates the module and local imports.
