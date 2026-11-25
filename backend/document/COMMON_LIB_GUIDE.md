# `core/common-lib.git/v2` Guide

`gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2` is used
by 7 of the 8 reviewed services (push-notification, sms-service,
email-service, email-gateway, as400-adapter, face-verification-service,
ekyc-service). `backend-for-frontend` is the one exception — it predates
wide adoption of this module and still carries equivalent code locally.

**Rule for new services: always depend on `core/common-lib.git/v2` for the
packages below. Do not re-implement them locally, and do not copy BFF's
local equivalents.** BFF is the outlier here, not the pattern to follow —
see ARCHITECTURE.md's divergence table.

## What's inside

| Package | Gives you | Don't do instead |
|---|---|---|
| `logger` | zap-backed structured logger, pre-configured field shape (service_id, request_id, level) | calling `zap.NewProduction()` yourself |
| `errorcodes` | the shared catalog of platform error codes | inventing per-service error code strings |
| `stderr` | error type that carries an `errorcodes` code + wraps the underlying cause + stacktrace | `cockroachdb/errors`, `fmt.Errorf` for anything crossing a package boundary |
| `stdresp` | the NEO response envelope builder (`status` / `service_id` / `message` / `data`) | hand-rolling `gin.H{...}` response maps |
| `ctxutils` | request-id / trace-id get/set on `context.Context` | passing IDs as separate function params |
| `middleware` | Gin middleware: request-id injection, access logging, panic recovery | writing your own recovery middleware |

(`internal/pool` and `internal/stacktrace` inside the module are private
implementation detail — don't import them directly; they're consumed through
`stderr`.)

## What's *not* inside (keep these local, in your service's own `pkg/`)

Things common-lib doesn't attempt to standardize, because they're
infra-specific rather than cross-cutting:

- HTTP client wrapper for calling other services (`pkg/client`)
- Kafka producer/consumer wiring (`pkg/kafka`)
- Redis client wiring (`pkg/redisclient`)
- Serialization helpers beyond what `encoding/json` gives you (`pkg/serializer`)
- Domain-specific utils: phone/bank-account formatting, masking, AES, etc. (`pkg/utils/*`)

BFF's `pkg/` is the reference for what this category of code looks like —
copy the *pattern*, not the dependency-free-of-common-lib choice.

## Using it in a new service

```go
import (
    "gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2/logger"
    "gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2/errorcodes"
    "gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2/stderr"
    "gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2/stdresp"
    "gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2/middleware"
)
```

Add it with the same version 7 other services already use unless you have a
specific reason to bump:

```bash
go get gitlab.awesome-poc-th.com/libero-engineering/core/common-lib.git/v2@v2.3.1
make vendor
```

See `template-service/go.mod` and `template-service/internal/middleware/`
for the wiring example.
