# My Tech Stack

Personal architecture + tooling blueprints, organized by stack area.
The backend blueprint is adapted from `tcr-tech-stack` (my work Go-service
standard, converged from 8 services); the frontend stacks are documented here
as well.

| Area | Contents |
|---|---|
| [`backend/`](backend/) | Go backend service blueprint — architecture, tech stack, common-lib usage, mocking strategy, CI/CD, scaffold-able `template-service/` |
| [`mobile/`](mobile/) | Mobile frontend blueprint — Swift / SwiftUI (iOS) and Kotlin / Jetpack Compose (Android) |
| [`web/`](web/) | Web frontend blueprint — Vite + bun + Vue 3 + TypeScript + shadcn UI + Tailwind CSS |

`.kiro/` at the repo root holds the Kiro steering and agent config that
drives the `backend/` scaffolding questionnaire — see
[backend/README.md](backend/README.md) for how to use it.

## Provenance / adaptation notes

- The backend area is copied-and-adapted from `~/Github/tcr-tech-stack` and
  is unchanged in substance. Its references to internal hosts
  (`gitlab.awesome-poc-th.com`, `core/common-lib.git/v2`, the NEO envelope)
  are that work platform's conventions — keep them when scaffolding work
  services, genericise them when reusing a skeleton elsewhere.
- [`mobile/`](mobile/document/TECH_STACK.md) and [`web/`](web/document/TECH_STACK.md)
  are the personal-stack choices above, written in the same blueprint style
  as the backend docs.
