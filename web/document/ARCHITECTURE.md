# Web Frontend Architecture Blueprint

Vue-idiomatic layering with the same intent as the backend blueprint: pages
are thin, logic lives in composables, all external I/O sits behind one
typed seam.

```
pages/ (route views)     — assemble composables + components; no API calls inline
components/
  ui/                    — shadcn-vue generated, vendored, owned here
  base/                  — our own small generic primitives
composables/             — reusable logic units (useXxx): state + side effects, no template
lib/api/                 — the only place fetch happens: endpoints, envelope decode, error mapping
stores/ (if Pinia)       — cross-view server/app state
types/                   — shared DTO types, mirroring backend models/ verbatim
assets/, lib/utils/      — static assets, pure helpers
```

## Data flow

```
route change → router guard (auth/permission) → page component
  → composable calls lib/api endpoint fn
  → api wrapper hits backend, unwraps the envelope (status/service_id/message/data)
  → typed DTO returned to composable → reactive state → page re-renders
  → error branch: backend errorcodes mapped to typed ApiError cases,
    surfaced via toast/inline — never string-matched on message
```

## Rules that make it scale

1. **One HTTP seam.** Components and pages never call `fetch` directly;
   `lib/api/<domain>.ts` exposes typed functions that return `Result<Dto>`.
2. **Composables own logic; components own markup.** Anything with `ref` +
   `watch` + async work graduates out of components into `useXxx`.
3. **Types mirror the backend.** DTO interfaces are generated or maintained
   to match backend `models/` field-for-field — the same contract is shared
   by the mobile blueprint, so drift is visible in one PR.
4. **Design only through tokens.** Tailwind utilities + shadcn CSS variables;
   no raw hex in components, no global CSS overrides outside `assets/`.
5. **Routes lazy-load by feature.** `components/` per feature directory if a
   view is non-trivial (`src/features/<name>/components`), mirroring the
   backend's domain-package rule: one capability = one directory that owns
   its pieces.

## Environment & config

- `VITE_API_BASE_URL` (and friends) via `import.meta.env`, with a
  `src/lib/env.ts` that validates presence at boot and fails loudly — the
  web equivalent of the backend's `CONFIG_FILE` discipline.
- `.env.example` committed; real values never.
