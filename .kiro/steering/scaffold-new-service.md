---
title: Scaffold a new BE service
description: Entry point for walking a user through creating a new Go backend service from this repo's blueprint.
---

# Scaffolding a new service

Trigger: the user asks to scaffold, create, bootstrap, or start a new
backend service in this workspace (phrases like "scaffold a new service",
"create a new BE service", "start a new microservice").

All blueprint content this flow reads and writes lives under `backend/`
(e.g. `backend/document/SCAFFOLD_CHECKLIST.md`, `backend/script/new-service.sh`,
`backend/template-service/`) — this `.kiro/` directory itself stays at the
workspace root, since that's where Kiro looks for steering/agents.

This is an **index**. Walk the sub-files below in order, in the same
conversation, asking the user before moving to the next one. Each sub-file
is a single stage — don't skip ahead or combine stages, and don't re-ask
something a later stage's spec file already fixes (e.g. don't re-derive the
infra checklist — it's `backend/document/SCAFFOLD_CHECKLIST.md`, read it
once and reuse it).

| Stage | File | Produces |
|---|---|---|
| 1 | [scaffold/01-naming-module.md](scaffold/01-naming-module.md) | service name, module path, port |
| 2 | [scaffold/02-infra-checklist.md](scaffold/02-infra-checklist.md) | cache/db/kafka/storage/firebase/mockoon/mocktcp answers |
| 3 | [scaffold/03-bruno-mockoon-setup.md](scaffold/03-bruno-mockoon-setup.md) | decides whether/how Bruno + Mockoon wiring happens (mostly deferred to skills, post-generation) |
| 4 | [scaffold/04-cicd-setup.md](scaffold/04-cicd-setup.md) | confirms which CI/CD stages to keep |
| 5 | [scaffold/05-generate.md](scaffold/05-generate.md) | recaps all answers, runs `backend/script/new-service.sh`, reports next steps |

Each stage ends by handing its answers forward — stage 5 is the only one
that actually writes anything to disk (by invoking the generator); stages
1-4 are pure Q&A that accumulate the flags stage 5 passes through.

If the user already knows all their answers up front and would rather skip
the back-and-forth, that's fine — collect everything in one pass instead of
five, then go straight to stage 5. The sub-files exist to structure the
conversation when the user doesn't already have answers ready, not to force
five round-trips on someone who does.
