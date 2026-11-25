# CI/CD Guide

All 8 services run on GitLab CI with the same pipeline shape (BFF's is the
most complete at 556 lines; as400-adapter's is the leanest at 216 — the
difference is almost entirely in how many test/integration stages are
present, not in structure). The template at
[`template-service/.gitlab-ci.yml`](../template-service/.gitlab-ci.yml) is
the trimmed, generic version of BFF's pipeline — delete stages your service
doesn't need rather than adding stages a bare service doesn't have yet.

## Stage pipeline

```
pre-build → build → unit-test → [docker-integration-test] → build-docker → notify
```

| Stage | Does | Runs when |
|---|---|---|
| `pre-build` | `go mod download` (with retry) + `go mod vendor`, using `GOPRIVATE`/`GOPROXY` so private modules (`common-lib.git/v2`, etc.) resolve; uploads `vendor/` as an artifact for later stages | every pipeline |
| `build` | compiles the binary with `-mod=vendor`, consuming the `vendor/` artifact from `pre-build` | every pipeline |
| `unit-test` | `go test ./...` with coverage | every pipeline |
| `docker-integration-test` | `docker compose -f docker-compose.integration.yaml up`, runs Bruno against it (see MOCKING_GUIDE.md) | services with integration tests wired up |
| `build-docker` | builds and pushes the image to ECR | `main` branch and tags only |
| `notify` | posts pipeline result (success/failure) to chat | always, as the last stage |

## Trigger rules

Pipelines run only for: an open merge request, a push to `main`, or a tag.
Plain commits on a feature branch with no open MR don't trigger a pipeline —
this avoids burning CI minutes on WIP branches nobody's reviewing yet.
`workflow.auto_cancel.on_new_commit: interruptible` cancels a ref's own
in-flight pipeline when a newer commit lands on that same ref (an MR's
previous run, or main's previous run) — it does not touch other refs'
pipelines.

## Private module resolution (the part that's easy to get wrong)

Every service depends on at least `core/common-lib.git/v2`, which lives on
the same private GitLab instance the pipeline runs on. `pre-build` sets this
up before `go mod download`:

```yaml
before_script:
  - echo "machine gitlab.awesome-poc-th.com login gitlab-ci-token password ${CI_JOB_TOKEN}" > ~/.netrc
  - chmod 600 ~/.netrc
  - git config --global url."https://gitlab.awesome-poc-th.com/".insteadOf "ssh://git@gitlab.awesome-poc-th.com/"
  - git config --global url."https://gitlab.awesome-poc-th.com/".insteadOf "git@gitlab.awesome-poc-th.com:"
  - go env -w GOPRIVATE=gitlab.awesome-poc-th.com
  - go env -w GOPROXY=http://goproxy-gitlab.awesome-poc-th.com:3000,direct
```

`vendor/` is then committed-via-artifact through the rest of the pipeline
(not re-fetched), which is also why local `make docker-build` depends on
`make vendor` first — see TECH_STACK.md's "Why vendoring."

## Go module / build cache

Pipelines share a runner-persisted cache keyed on `go.sum`, holding
`GOMODCACHE` and `GOCACHE`. Bump the cache key prefix (e.g. `1_go` → `2_go`)
if the cache ever gets corrupted and needs a hard reset — don't try to
selectively clear it.

## Worked example — adding a new stage-gated job

Say your new service needs a `lint` job before `build`. Add it to the same
stage group, gated the same way as `build`:

```yaml
lint:
  stage: pre-build          # runs alongside pre-build, not a new stage
  tags: ["linux"]
  image: public.ecr.aws/docker/library/golang:${GO_VERSION}
  needs:
    - job: pre-build
      artifacts: true
  script:
    - go run github.com/golangci/golangci-lint/cmd/golangci-lint@v2.1.0 run
  rules:
    - when: on_success
```

Don't invent a new top-level `stage:` for a single job unless it genuinely
needs to run after some stages and before others — most new jobs belong
inside an existing stage.

## What to trim for a simple/leaf service

A service with no outbound HTTP calls and no integration test stack yet
(the sms-service/email-service shape) can drop `docker-integration-test`
entirely and keep `pre-build → build → unit-test → build-docker → notify`.
Add `docker-integration-test` back in once `skill::bruno-integration-test`
has set up the IT stack.
