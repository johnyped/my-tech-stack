# Repository guidance

- Use `backend/` for the Go service blueprint and generator. The `web/` and
  `mobile/` directories contain stack and architecture guides.
- Use `backend/script/new-service.sh` as the generation entry point and
  `backend/script/new_service.py` as the implementation. Keep the Makefile,
  Codex skill and agent, and Kiro questionnaire aligned with the generator.
- Use the [scaffold checklist](backend/document/SCAFFOLD_CHECKLIST.md) as the
  reference for supported options. For interactive scaffolding, use
  `$scaffold-new-service` and collect only the inputs still needed.
- Generate test services in temporary directories. Before using `--force`,
  check that the target is a marked scaffold and the user has authorized
  replacing its contents.
- Keep service name and port placeholders consistent across Go config,
  Dockerfile, Makefile, Compose, and Bruno. Dependency ports belong to the
  dependency configuration.
- Verify backend scaffold changes from the repo root with
  `python3 -B -m unittest discover -s backend/script/tests -v`.
  Docker Compose is required for the Compose checks. A full Go build and
  liveness check require access to the private work-platform module.
