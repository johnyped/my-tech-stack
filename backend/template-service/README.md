# TODO-service-name

Scaffolded from this repo's blueprint. Replace this README with the real
service's once you're done wiring it up.

## If you copy/pasted this folder by hand

Find-and-replace these, in this order:

1. Replace the full placeholder Go module path everywhere it appears,
   including imports in Go files. See `document/TECH_STACK.md` for module naming.
2. Replace the remaining service-name placeholder with your kebab-case name
   and the service-port placeholder with your port (normally `8080`).
3. Delete whichever of these you don't need, per your checklist answers
   (`document/SCAFFOLD_CHECKLIST.md`):
   - `internal/storage/` (no `--storage s3`)
   - `mocktcp/` (no `--mocktcp yes`)
   - `test/integration/mockoon/`, `test/integration/scenarios.yaml` (no `--mockoon yes`)
   - the matching commented blocks in `docker-compose.yaml` / `config/*.yaml`
4. Run `make setup && make config-init && make run` and confirm
   `curl localhost:TODO-service-port/health/liveness` returns `200`.
5. Once you have real routes, load `skill::bruno-service-docs` to generate
   `bruno/` for real, and `skill::bruno-integration-test` to replace the
   placeholder scripts in `script/integration/`.

## If you used `make new-service` or the Kiro questionnaire instead

Steps 1-3 above are already done for you — start at step 4.
