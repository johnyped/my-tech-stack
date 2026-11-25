# My Tech Stack

Personal architecture + tooling blueprints, organized by stack area.
The backend blueprint is adapted from `tcr-tech-stack` (my work Go-service
standard, converged from 8 services); the frontend stacks are documented here
as well.

| Area | Contents |
|---|---|
| [`backend/`](backend/) | Go backend service blueprint — architecture, tech stack, request middleware, mocking strategy, CI/CD, scaffold-able `template-service/` |
| [`mobile/`](mobile/) | Mobile frontend blueprint — Swift / SwiftUI (iOS) and Kotlin / Jetpack Compose (Android) |
| [`web/`](web/) | Web frontend blueprint — Vite + bun + Vue 3 + TypeScript + shadcn UI + Tailwind CSS |

## Scaffold a Go service

From the repository root, run:

```sh
backend/script/new-service.sh \
  --name my-service \
  --module example.com/my-service
```

The generator creates a sibling directory by default. Use `--out` to choose a
different destination and `--port` to override the default port of `8080`.
Use `--force` only when replacing a marked scaffold whose contents you intend
to discard. See the [backend guide](backend/README.md) and
[scaffold checklist](backend/document/SCAFFOLD_CHECKLIST.md) for all options.

With Codex, invoke `$scaffold-new-service` to walk through the questionnaire.
To delegate generation after agreeing on the answers and destination, ask
Codex to use the `scaffold-new-service` custom agent. See the
[skill](.agents/skills/scaffold-new-service/SKILL.md) and
[agent definition](.codex/agents/scaffold-new-service.toml). In Kiro, ask
“scaffold a new service” to start the [same questionnaire](.kiro/steering/scaffold-new-service.md).

## Verify the scaffold

```sh
python3 -B -m unittest discover -s backend/script/tests -v
```

Install Docker Compose to run the configuration checks for the minimal, Redis,
and full-option Valkey fixtures. The [Backend scaffold workflow](.github/workflows/backend-scaffold.yml)
runs these checks on backend changes. The [manual runtime job](backend/README.md#scaffold-verification)
builds and starts a generated service using public Go modules.

## Provenance / adaptation notes

- The backend blueprint was adapted from `tcr-tech-stack`. The scaffold uses
  public Go modules and local request middleware.
- [`mobile/`](mobile/document/TECH_STACK.md) and [`web/`](web/document/TECH_STACK.md)
  are the personal-stack choices above, written in the same blueprint style
  as the backend docs.
