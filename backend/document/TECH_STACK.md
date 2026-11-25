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
| Structured logging | `go.uber.org/zap`, wrapped by `common-lib.git/v2/logger` | don't call zap directly in handlers — use the common-lib wrapper so log shape is consistent across services |
| Request validation | `go-playground/validator/v10` | struct tag validation on request-bound structs |
| Error types | `core/common-lib.git/v2/errorcodes` + `stderr` | not `cockroachdb/errors` (BFF-only, not adopted — see ARCHITECTURE.md) |
| Response envelope | `core/common-lib.git/v2/stdresp` | the NEO envelope: `status` / `service_id` / `message` / `data` |
| Context propagation | `core/common-lib.git/v2/ctxutils` | request-id, trace-id helpers |
| Common middleware | `core/common-lib.git/v2/middleware` | request-id injection, recovery, access log |

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

Every service depends on at least one private module under
`gitlab.awesome-poc-th.com/...` (`core/common-lib.git/v2` at minimum). CI
runners and Docker builds don't carry the credentials to fetch private Go
modules at build time, so `vendor/` is committed and `go build -mod=vendor`
is used everywhere. `make vendor` must run (and succeed, with local
credentials) before `make docker-build`.

## Module naming — no fixed convention today, pick deliberately

The 8 services disagree on this:

- Plain name, no path: `backend-for-frontend`, `push-notification-service`, `ekyc-service`
- Wrong/placeholder path: `github.com/neo/sms-service`
- Real gitlab group path: `gitlab.awesome-poc-th.com/libero-engineering/core/email-gateway`, `.../core/as400adapter`, `.../auxiliary-systems/ekyc/face-verification-service`

For a new service, use the real gitlab group path
(`gitlab.awesome-poc-th.com/libero-engineering/<group>/<service-name>`) —
plain names and placeholder paths both break `go get` for anyone depending on
this service later. Pick `<group>` to match where the service actually lives
in GitLab (`core` for platform-wide services, `auxiliary-systems/<system>`
for a system-specific one, by precedent from face-verification-service and
ekyc-service). The scaffold generator's `--module` flag requires this as
input rather than defaulting it — see SCAFFOLD_CHECKLIST.md.
