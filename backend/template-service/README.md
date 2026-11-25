# TODO-service-name

Scaffolded from this repo's blueprint. Replace this README with the real
service's once you're done wiring it up.

## If you copy/pasted this folder by hand

Find-and-replace these, in this order:

1. `TODO-service-name` → your service's kebab-case name, everywhere (go.mod
   module basename, Dockerfile binary name, docker-compose container name,
   `.gitlab-ci.yml` `SERVICE_NAME`, bruno env URLs, `.golangci.yml`
   local-prefixes).
2. `gitlab.awesome-poc-th.com/libero-engineering/TODO-group/TODO-service-name`
   in `go.mod` → your real module path (see `document/TECH_STACK.md`
   "Module naming" for how to pick `TODO-group`).
3. Delete whichever of these you don't need, per your checklist answers
   (`document/SCAFFOLD_CHECKLIST.md`):
   - `internal/storage/` (no `--storage s3`)
   - `mocktcp/` (no `--mocktcp yes`)
   - `test/integration/mockoon/`, `test/integration/scenarios.yaml` (no `--mockoon yes`)
   - the matching commented blocks in `docker-compose.yaml` / `config/*.yaml`
4. Run `make setup && make config-init && make run` and confirm
   `curl localhost:8080/health/liveness` returns `200`.
5. Once you have real routes, load `skill::bruno-service-docs` to generate
   `bruno/` for real, and `skill::bruno-integration-test` to replace the
   placeholder scripts in `script/integration/`.

## If you used `make new-service` or the Kiro questionnaire instead

Steps 1-3 above are already done for you — start at step 4.
