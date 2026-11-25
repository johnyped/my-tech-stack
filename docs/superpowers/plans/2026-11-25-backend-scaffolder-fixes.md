# Backend scaffolder fixes

Fix output safety and make every generated service use consistent names, ports, and dependency settings. Keep the existing CLI entry points and work-platform defaults. The portable preset can be a separate follow-up if personal reuse warrants it.

## 1. Protect the output directory

- [x] Add Python standard-library regression tests under `backend/script/tests/`. Use temporary directories and cover invalid names, invalid ports, collisions, and a simulated failure before publishing the generated service.
- [x] In `backend/script/new_service.py`, validate kebab-case names and integer ports in the range 1–65535 before any writes. Check the resolved output against the blueprint checkout, its ancestors, and symlink targets.
- [x] Stage the scaffold in a temporary directory beside the destination. Publish a new scaffold only after rendering succeeds. For `--force`, require a generator-owned marker, keep the old output until the new scaffold is ready, and restore it if publishing fails. Document that older unmarked output requires a fresh destination.
- [x] Update `backend/document/SCAFFOLD_CHECKLIST.md` and the Codex scaffold skill and agent to reflect the new overwrite behavior. Check the Kiro guidance for the same assumption.

**Acceptance.** Normal generation works in a new directory. An unrelated directory keeps its contents even with `--force`. A failed render keeps the previous scaffold intact.

## 2. Render service names and ports consistently

- [x] Add fixtures for the default port and for port `9191` with Kafka enabled. Assert the service port and Kafka UI container port independently.
- [x] In `backend/script/new_service.py`, explicitly render extensionless template files including `Makefile` and `Dockerfile`. Replace the blanket `8080` rewrite with a dedicated service-port token. Audit application config, Compose, Dockerfile, Makefile, Bruno environments, and service documentation for service-port references. Keep the documented copy-and-edit path usable and verify that no render tokens remain in generated output.

**Acceptance.** The configured port, health checks, Docker mapping, and Bruno local URLs agree. The generated binary and image use the requested service name. Kafka UI still maps to container port `8080` when the service uses `9191`.

## 3. Make optional dependencies match their flags

- [x] Generate Redis and Valkey separately and validate the selected image, startup command, and health check. Preserve the cache address expected by the service config where possible.
- [x] Make the starter Mockoon command point to the shipped fixture and verify it starts with that fixture. Keep scenario switching in the separate integration-test workflow.
- [x] Put MinIO on the same Compose network as the app and verify the configured S3 hostname resolves from that network.

**Acceptance.** Each flag produces the intended Compose service, with a resolvable fixture path or network address where applicable. Redis and Valkey each have coverage.

## 4. Run scaffold smoke tests in CI

- [x] Run `python3 -m unittest discover -s backend/script/tests` on the review pipeline. Generate minimal, Redis, and custom-port full-option Valkey services in temporary directories. Check rendered values, enabled services, and files referenced by Compose, then run `docker compose config --quiet` for each fixture.
- [x] Add a separate runtime check on a runner with access to the private Go module. Run the generated service's `make vendor`, build the container, start it, and require HTTP 200 from `/health/liveness`. Exercise the selected dependency containers in the full-option fixture and tear the stack down afterward.

**Acceptance.** Static checks run without private module access. The authenticated runtime job verifies boot and health, and reports dependency or startup failures separately from rendering failures.

## Verification status

Eight local regression tests passed, including Compose validation for minimal,
Redis, and full-option Valkey output. The full runtime acceptance check is
pending on an authenticated runner. Local vendoring stalled while fetching the
private common-lib module, whose HTTPS Git endpoint requires authentication.
