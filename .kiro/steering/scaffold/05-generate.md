---
title: Scaffold stage 5 — generate
description: Final stage of the new-service scaffold flow — recaps answers, runs the generator, reports next steps.
---

# Stage 5 — Generate

1. Recap every answer collected in stages 1-4 as a single summary (name,
   module, port, cache, db, kafka, storage, firebase, mockoon, mocktcp) and
   ask for a final go-ahead before running anything — this is a disk-write
   action (creates a new directory), not a read.
2. On confirmation, run:
   ```
   backend/script/new-service.sh --name <name> --module <module> --port <port> \
     --cache <cache> --db <db> --kafka <kafka> --storage <storage> \
     --firebase <firebase> --mockoon <mockoon> --mocktcp <mocktcp>
   ```
   (omit a flag entirely if its stage left it at the default `no` — the
   script defaults match stage 2's recommendations, so this is just
   passing through what was collected, not re-deciding anything.)
3. The script prints its own next steps (go get / make vendor / make run /
   curl health check) — relay that output, don't re-derive it.
4. Remind the user of the two follow-on skills from stage 3 once they have
   real routes: `skill::bruno-service-docs` and, if `--mockoon yes`,
   `skill::bruno-integration-test`.
5. If stage 4 flagged `docker-integration-test` as removable, remind them
   of that now that the file actually exists on disk to edit.

This is the only stage that touches the filesystem. To regenerate with revised
answers, use a new destination or get explicit authorization to replace a
marked scaffold with `--force`. Unmarked directories and Git repositories
require a fresh destination.
