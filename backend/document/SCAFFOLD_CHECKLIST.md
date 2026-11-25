# New Service Scaffold Checklist

This is the single spec for "what questions does a new service need
answered before it can be generated." `script/new-service.sh` (flags),
`.kiro/steering/scaffold/02-infra-checklist.md` (conversational), and
`.kiro/agents/scaffold-new-service.json` all implement *this* list — if you
change a question here, update the flag/prompt in all three.

## Identity (always asked, no default)

| Question | Flag | Notes |
|---|---|---|
| Service name | `--name` | kebab-case, becomes the binary name and container name |
| Module path | `--module` | full Go module path (for example, `example.com/team/my-service`) — see TECH_STACK.md's module-naming note; never defaulted |
| Port | `--port` | default `8080` if not given |

## Infra checklist (each is yes/no, default no)

| Question | Flag | If yes, adds |
|---|---|---|
| Cache / session store? | `--cache redis\|valkey\|no` | `redis`/`valkey` compose service, `pkg/redisclient` wiring, config block |
| Relational DB? | `--db postgres\|no` | `postgres` compose service, `pgx/v5` + `sqlc` + `db/migrations`, config block |
| Event streaming? | `--kafka yes\|no` | `kafka` + `kafka-init` + `kafka-ui` compose services, `pkg/kafka` wiring |
| Object storage? | `--storage s3\|no` | `minio` compose service (dev), AWS S3 adapter interface + stub |
| Push notification (Firebase)? | `--firebase yes\|no` | Firebase adapter interface + stub |
| Calls another HTTP service? | `--mockoon yes\|no` | Mockoon compose service + `scenarios.yaml` + `bruno/Integration Tests/` wiring. Follow the decision rule in MOCKING_GUIDE.md §2 — don't default this to yes just because it's available |
| Non-HTTP downstream (e.g. raw TCP)? | `--mocktcp yes\|no` | `mocktcp/` stub server, matching your service's framing |

Two or more "yes" answers compose independently — e.g. `--cache redis --db
postgres --kafka yes` produces a compose file with all three services plus
the app, same as BFF's actual stack.

## What never gets asked (fixed by the blueprint, not a choice)

- Local request ID middleware and a Zap logger (see REQUEST_UTILITIES.md)
- Health endpoints (`/health/liveness`, `/health/readiness`, `/ping`) — always present
- `golangci-lint` v2 config — always included
- CI/CD stage skeleton — always included; trim stages per CICD_GUIDE.md's "what to trim" note if the service turns out to be a leaf service
- Bruno `Health/` + skeleton — always present; the rest of `bruno/` comes from `skill::bruno-service-docs` once routes exist

## Example runs

```bash
# Leaf service, no infra: sms-service shape
script/new-service.sh --name my-sms-relay \
  --module example.com/team/my-sms-relay

# BFF-shape: cache + kafka + storage, calls other HTTP services
script/new-service.sh --name my-bff \
  --module example.com/team/my-bff \
  --cache redis --kafka yes --storage s3 --mockoon yes

# ekyc-shape: DB + cache, calls other HTTP services
script/new-service.sh --name my-kyc-flow \
  --module example.com/identity/my-kyc-flow \
  --db postgres --cache redis --mockoon yes
```

### Regenerating a service

Use a fresh destination for existing projects and older, unmarked output.
With `--force`, the generator replaces a marked scaffold after the new output
has been rendered successfully. This replaces all local changes in that
directory. The generator rejects Git repositories as replacement targets.
