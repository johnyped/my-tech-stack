---
title: Scaffold stage 1 — naming & module path
description: First stage of the new-service scaffold flow — collects name, module path, port.
---

# Stage 1 — Naming & module path

Ask these in order. For each, state the recommendation before asking, same
as the rest of this blueprint's decisions were made.

1. **Service name** (kebab-case). Becomes the binary name, Docker container
   name, Bruno collection base URL label. No default — always ask.

2. **GitLab group**, to build the module path
   `gitlab.awesome-poc-th.com/libero-engineering/<group>/<service-name>`.
   See `backend/document/TECH_STACK.md` "Module naming" for why there's no single
   fixed default. Recommend `core` unless the service clearly belongs to a
   specific existing system (e.g. `auxiliary-systems/ekyc` the way
   face-verification-service and ekyc-service do) — ask the user which
   fits, don't guess silently.

3. **Port**. Recommend `8080` (what all 8 reviewed services use) unless the
   user has a reason to deviate.

Hold these three answers — service name, module path, port — and carry
them into stage 2. Don't invoke the generator yet.
