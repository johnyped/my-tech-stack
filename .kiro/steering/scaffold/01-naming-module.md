---
title: Scaffold stage 1 — naming & module path
description: First stage of the new-service scaffold flow — collects name, module path, port.
---

# Stage 1 — Naming & module path

Ask these in order. For each, state the recommendation before asking, same
as the rest of this blueprint's decisions were made.

1. **Service name** (kebab-case). Becomes the binary name, Docker container
   name, Bruno collection base URL label. No default — always ask.

2. **Go module path**. Use the path for the repository where the service will
   live. For a local example, suggest `example.com/<service-name>`.

3. **Port**. Recommend `8080` (what all 8 reviewed services use) unless the
   user has a reason to deviate.

Hold these three answers — service name, module path, port — and carry
them into stage 2. Don't invoke the generator yet.
