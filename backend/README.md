# BE Service Blueprint (adapted from tcr-tech-stack)

Unified architecture, tech stack, and project-setup template for Go backend
services on this platform, derived from a review of 8 existing services:

`backend-for-frontend` (reference/master), `push-notification`, `sms-service`,
`email-service`, `email-gateway`, `as400-adapter`, `face-verification-service`,
`ekyc-service`.

This is not 8 separate reviews. It's one converged standard — where the 8
services agreed, that's the standard; where `backend-for-frontend` (BFF)
disagreed with the rest, the majority convention wins for anything
cross-cutting (see [COMMON_LIB_GUIDE.md](document/COMMON_LIB_GUIDE.md)).

## Start here

| I want to... | Go to |
|---|---|
| Understand the layering / request flow / adapter pattern | [document/ARCHITECTURE.md](document/ARCHITECTURE.md) |
| See the full tech stack (frameworks, libs, infra) and why | [document/TECH_STACK.md](document/TECH_STACK.md) |
| Know what `core/common-lib.git/v2` gives me vs. local `pkg/` | [document/COMMON_LIB_GUIDE.md](document/COMMON_LIB_GUIDE.md) |
| Mock a downstream dependency in tests | [document/MOCKING_GUIDE.md](document/MOCKING_GUIDE.md) |
| Set up or understand the CI/CD pipeline | [document/CICD_GUIDE.md](document/CICD_GUIDE.md) |
| See the new-service questionnaire as a spec | [document/SCAFFOLD_CHECKLIST.md](document/SCAFFOLD_CHECKLIST.md) |
| Start a brand-new service | see **Scaffolding a new service** below |

## Scaffolding a new service

Three equivalent paths — pick whichever fits how you work. All three end up
calling the same generator, so none of them can drift from the others.

1. **Copy/paste** — copy [`template-service/`](template-service/) to your new
   repo, rename the module, follow the `TODO`s. No tooling required.
2. **Scripted** — run the generator directly with your answers as flags
   (or `make new-service NAME=... MODULE=... CACHE=redis ...` from this
   repo root — both call the same `script/new-service.sh`):
   ```bash
   script/new-service.sh \
     --name my-new-service \
     --module gitlab.awesome-poc-th.com/libero-engineering/core/my-new-service \
     --cache redis --db postgres --kafka no --storage no --mockoon no --mocktcp no --firebase no
   ```
3. **Kiro-driven** — open this workspace in Kiro and say "scaffold a new
   service." Kiro reads [../.kiro/steering/scaffold-new-service.md](../.kiro/steering/scaffold-new-service.md)
   and its sub-files, asks you the checklist one question at a time, then
   calls the same `script/new-service.sh` from option 2. Or invoke the
   dedicated agent directly: `../.kiro/agents/scaffold-new-service.json`.

## Repo layout

```
document/            architecture, tech stack, common-lib, mocking, CI/CD guides
template-service/    copy/paste-able skeleton (the actual scaffold)
script/new-service.sh   generator the scaffold flow and Makefile target both call
.kiro/steering/      Kiro-driven questionnaire (ad hoc, any session)
.kiro/agents/        dedicated scaffolding agent
```

## Scaffold verification

From the repository root, run `python3 -B -m unittest discover -s backend/script/tests -v`.
Docker Compose enables the generated configuration checks. GitHub Actions runs
these checks on backend changes without private module access.

The Backend scaffold workflow also exposes a manual runtime job on the default
branch. It needs a Linux runner labeled `backend-scaffold` with Python, Go
1.26, Make, curl, Docker Compose, and configured access to the private GitLab
Go module. Run it on a dedicated runner so the fixture's ports are available.
