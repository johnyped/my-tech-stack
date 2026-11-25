---
title: Scaffold stage 2 — infra checklist
description: Second stage of the new-service scaffold flow — walks the infra checklist from document/SCAFFOLD_CHECKLIST.md.
---

# Stage 2 — Infra checklist

Read `backend/document/SCAFFOLD_CHECKLIST.md` first — it's the single spec for
these questions; this file only tells you how to *ask* them, not what they
are (don't let the two drift apart by re-deriving the list here).

Ask each checklist item **one at a time**, in the table's order, with your
recommended default stated before asking:

1. Cache / session store — default recommendation: **no**, unless the
   service needs session state, rate-limiting, or a distributed lock (most
   leaf services, like sms-service/email-service's shape, don't).
2. Relational DB — default: **no**, same reasoning; only "yes" if this
   service is the owner of some persisted record (like push-notification,
   email-service, ekyc-service).
3. Event streaming (Kafka) — default: **no**. Only BFF, among the 8
   reviewed services, produces to Kafka.
4. Object storage (S3) — default: **no**. Only needed if the service reads
   or writes files/blobs.
5. Firebase (push) — default: **no**, unless this service sends push
   notifications directly (push-notification's job, not most services').
6. Calls another HTTP service (feeds the Mockoon decision in stage 3) —
   default: **no**. If yes, don't decide Mockoon yet — that's stage 3's job,
   this question's answer is just input to it.
7. Non-HTTP downstream (raw TCP etc.) — default: **no**. Rare — only
   as400-adapter among the 8 needs this.

Carry all seven answers into stage 3.
