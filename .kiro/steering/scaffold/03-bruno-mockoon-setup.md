---
title: Scaffold stage 3 — Bruno & Mockoon
description: Third stage of the new-service scaffold flow — decides the Mockoon answer and sets expectations for Bruno docs.
---

# Stage 3 — Bruno & Mockoon

Read `backend/document/MOCKING_GUIDE.md` §2 first — this stage applies its decision
rule, it doesn't restate it.

1. Using stage 2's answer to "calls another HTTP service": if **no**, set
   `--mockoon no` and move on without asking anything else — don't offer
   Mockoon to a service that has no outbound HTTP calls.
2. If **yes**, ask: "across your test scenarios, do you need more than one
   response shape per downstream (happy path, not-found, downstream error,
   etc.), or is happy-path enough?" If happy-path is enough, recommend
   `--mockoon no` (plain Bruno tests against the service itself cover it).
   Otherwise `--mockoon yes`.
3. Using stage 2's answer to "non-HTTP downstream": pass straight through
   as `--mocktcp yes|no` — no further question needed, Mockoon doesn't
   apply to that kind at all.

**Set expectations, don't generate Bruno content yet:** the generator
(stage 5) only produces the static skeleton (`Health/`, `Example Domain/`,
`environments/`) — tell the user that once they have real routes, the next
step is loading `skill::bruno-service-docs` (API docs) and, if
`--mockoon yes`, `skill::bruno-integration-test` (IT stack + Mockoon
scenarios + runner scripts). Don't attempt to write Bruno YAML or Mockoon
JSON yourself in this flow — that's what those two skills are for.

Carry `--mockoon` and `--mocktcp` into stage 4.
