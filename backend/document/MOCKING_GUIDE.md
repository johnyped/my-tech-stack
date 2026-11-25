# Mocking & Stubbing Guide

A new service can have up to four kinds of "fake the outside world" need.
They solve different problems — don't collapse them into one mechanism.

| Downstream kind | Example | Mechanism | When it runs |
|---|---|---|---|
| HTTP service, everyday local dev | calling auth-service while coding a BFF feature | in-process Go `mockservice/` package | `make run-docker`, zero extra containers |
| HTTP service, integration test scenarios | same call, but you need 5 different response shapes for 5 test cases | Mockoon + `scenarios.yaml` | `make integration-test-*` |
| Non-HTTP downstream | AS400 TCP link | custom mock TCP server (`mocktcp/`) | both local dev and integration tests |
| SDK-based external | AWS S3, Kafka, Firebase | adapter interface + stub implementation | unit tests, and local dev via `STUB_MODE` |

## 1. In-process Go mock (local dev convenience)

Pattern from BFF's `mockservice/`: one `.go` file per downstream, a
`registry.go` that wires them into a local `httptest`-style server or a
second Gin instance started alongside the app in `docker-compose.yaml`. Use
this when you just need the app to boot and respond with *something*
plausible while you're working on a feature — not for asserting specific
test scenarios.

Only build this if your service calls other HTTP services in the first
place. Most of the 8 services don't (sms-service, email-service,
email-gateway are leaf services with no outbound HTTP calls) — skip this
entirely if that's you.

## 2. Mockoon + scenarios (integration test scenario matrix)

Pattern from ekyc-service. Use when Bruno integration tests need the
downstream to return *different* payloads for different test cases (happy
path, not-found, downstream-5xx, slow response, etc.) — Mockoon's per-route
response rules and `scenarios.yaml`'s scenario switching are what make that
practical; the in-process Go mock from (1) isn't built for that.

**Decision rule (same one `skill::bruno-integration-test` uses) — don't add
Mockoon unless this is actually true for your service:**

1. Does your service make outbound HTTP calls to another service at all? If
   no, you don't need Mockoon — plain Bruno request/response tests against
   your own service are sufficient.
2. If yes, and you need more than one response shape per downstream across
   your test scenarios, use Mockoon. If every test case is happy-path only,
   you may still not need it.

**Don't hand-write the Mockoon/Bruno wiring.** Load and follow:

- `skill::bruno-service-docs` — generates the Bruno API-doc collection (`bruno/`) from your service's own routes.
- `skill::bruno-integration-test` — sets up the full IT stack: `bruno/Integration Tests/`, Mockoon datasources + `scenarios.yaml` (only if step 1 above says yes), runner scripts, and Makefile targets.

Both skills already encode every YAML schema rule, the scenario-numbering
convention, and the validator script. Re-deriving that here would just drift
out of sync with them.

## 3. Custom TCP mock (non-HTTP downstreams)

as400-adapter's downstream is a raw TCP socket with a fixed-width framed
message protocol (`internal/models/message`), not HTTP — Mockoon can't mock
it. The pattern is a small standalone TCP server that:

- listens on the same port shape as the real AS400 link,
- parses the same frame header (`tcp_header.go`-equivalent) and operation code,
- returns canned responses keyed by operation code, the same way Mockoon
  keys responses by HTTP route.

See `template-service/mocktcp/` for the stub. If your new service talks to
a non-HTTP downstream that isn't AS400, copy this pattern rather than
inventing a new mocking mechanism per protocol.

## 4. SDK-based external adapters (AWS, Kafka, Firebase, ...)

These aren't "downstream services" you call over a wire you control — they're
vendor SDKs. Mockoon and mocktcp don't apply. Instead:

1. Define an interface in `internal/<concern>/` for the operations you
   actually use (see ARCHITECTURE.md's adapter pattern example).
2. Real implementation wraps the SDK.
3. Stub implementation is an in-memory fake, selected via config/env
   (`STUB_MODE=true` or per-dependency env flags), wired once in
   `cmd/main.go`.
4. Unit tests depend on the interface and inject the stub directly — no env
   flag needed at that layer.

This is how you get AWS/Kafka/Firebase calls testable without hitting real
infrastructure in CI, without adding a protocol-specific mock server for
each one.
