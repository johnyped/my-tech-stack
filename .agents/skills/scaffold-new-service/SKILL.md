---
name: scaffold-new-service
description: Scaffold a new Go backend service from this repository's backend blueprint. Use when asked to scaffold, bootstrap, or start a new backend service. For existing services, use the relevant maintenance workflow.
---

# Scaffold a new Go backend service

Run the questionnaire in the main Codex conversation. Carry answers forward and ask only for missing choices, one at a time, with the recommended default first. If the user supplies all answers upfront, move directly to the generation recap. Use repository-root-relative paths below.

1. **Naming.** Collect a kebab-case service name, Go module path, and port. There is no default name. Use the full module path supplied by the user or derive it from the intended repository URL. For a local example, use `example.com/<service-name>`. See [module naming](../../../backend/document/TECH_STACK.md). Recommend port `8080`.
2. **Infrastructure.** Read the [scaffold checklist](../../../backend/document/SCAFFOLD_CHECKLIST.md). Ask in order about cache (`redis|valkey|no` for session state, rate limits, or locks), DB (`postgres|no` for owned records), Kafka (event streaming), S3 (files or blobs), Firebase (direct push), outbound HTTP, and non-HTTP downstreams such as raw TCP. Recommend `no` unless the service needs the capability.
3. **Bruno and mocks.** Apply §2 of the [mocking guide](../../../backend/document/MOCKING_GUIDE.md). With no outbound HTTP, use `--mockoon no`. Otherwise ask whether tests need multiple downstream response shapes. Use `--mockoon yes` for multiple shapes and recommend `no` for happy-path-only tests. Map non-HTTP downstreams directly to `--mocktcp yes|no`. The scaffold ships starter Bruno files; route-specific docs and scenarios belong to the follow-on `$bruno-service-docs` and, when Mockoon is enabled, `$bruno-integration-test` skills if installed.
4. **CI/CD.** Read “What to trim” in the [CI/CD guide](../../../backend/document/CICD_GUIDE.md). The template includes the full GitLab pipeline. For a leaf service without outbound HTTP or a Mockoon stack, flag `docker-integration-test` as a candidate for manual removal. This stage adds no flags.
5. **Generate.** Recap the name, module, port, infrastructure and mock selections, and output directory. Ask for the go-ahead if the user has yet to authorize generation. From the repository root run the existing generator with the agreed arguments. Omit flags whose values match the CLI defaults.

   ```text
   backend/script/new-service.sh --name <name> --module <module> --port <port> \
     --cache <cache> --db <db> --kafka <kafka> --storage <storage> \
     --firebase <firebase> --mockoon <mockoon> --mocktcp <mocktcp>
   ```

The default destination is a direct sibling of this repository named after the service. Use `--out` for an agreed destination within this workspace or a new direct sibling directory. Check the destination before generation. If it exists, report the collision and offer a fresh destination. The generator accepts `--force` only for a marked scaffold without a Git repository. Use it only with explicit authorization to replace that directory and its contents.

Verify that the generated directory and key files exist. Report the path and relay the generator's relevant next steps, including the CI/CD note and Bruno follow-ons where applicable. Leave application routes, Bruno scenarios, and pipeline edits to follow-up tasks.
