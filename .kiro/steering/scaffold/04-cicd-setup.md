---
title: Scaffold stage 4 — CI/CD
description: Fourth stage of the new-service scaffold flow — confirms which pipeline stages to keep.
---

# Stage 4 — CI/CD

Read `backend/document/CICD_GUIDE.md` first, specifically "What to trim."

The generator always includes the full `.gitlab-ci.yml` template
(`pre-build → build → unit-test → docker-integration-test → build-docker →
notify`) — this stage doesn't change what gets generated, it just sets
expectations and flags the one stage that's conditional:

1. Tell the user: every stage except `docker-integration-test` applies
   regardless of their checklist answers.
2. If stage 3 ended with `--mockoon no` and stage 2's "calls another HTTP
   service" was also no (i.e. this is a leaf service with no integration
   test stack), say so explicitly: recommend deleting the
   `docker-integration-test` job from the generated `.gitlab-ci.yml` after
   scaffolding, since it has nothing to exercise yet. Don't delete it
   automatically — the user may add Mockoon later and want it back; just
   flag it as the one piece of manual cleanup in this generated file.
3. No new flags come out of this stage — it's advisory only. Move to
   stage 5 with everything collected so far.
