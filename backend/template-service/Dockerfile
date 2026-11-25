# Build stage
FROM public.ecr.aws/docker/library/golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod files and the pre-vendored dependency tree (run `make vendor`
# locally/in CI before this — private modules can't be fetched from inside
# the build, see document/TECH_STACK.md "Why vendoring")
COPY go.mod go.sum ./
COPY vendor/ vendor/

COPY . .

# TODO: replace TODO-service-name with the real binary name (match Makefile's SERVICE_NAME)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -mod=vendor \
    -ldflags="-w -s" \
    -o TODO-service-name \
    cmd/main.go

# Final stage - minimal runtime image
FROM public.ecr.aws/docker/library/alpine:3.21

RUN apk --no-cache add ca-certificates tzdata wget

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

ENV CONFIG_FILE=/app/config/config.yaml

COPY --from=builder /app/TODO-service-name .
COPY --from=builder /app/config ./config

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health/liveness || exit 1

CMD ["./TODO-service-name"]
