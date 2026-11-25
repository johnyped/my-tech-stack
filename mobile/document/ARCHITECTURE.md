# Mobile Architecture Blueprint

The same three-layer shape on both platforms. Names differ (Swift/
SwiftUI vs Kotlin/Compose), responsibilities don't.

```
View (SwiftUI / Composable)
  renders State — holds no business logic
    ↑
Feature Model (@Observable "view model" / Compose ViewModel)
  owns screen state, orchestrates use cases, exposes Intent → State
    ↑
Domain layer
  use cases / repositories — plain code, no UI or platform imports
    ↑
Data layer
  API client (backend envelope), DTO mapping, local cache (URLCache/Codable
  on iOS; Room/DataStore on Android)
```

## Feature-first layout

```
app/                 entry point, root navigator, dependency wiring
features/<name>/     one module per user-facing capability
  view/              screens + small reusable view components
  model/             the observable view model
  domain/            use cases + repository protocol/interface
  data/              DTOs, API endpoints, cache data source
shared/              api client, auth, networking, design tokens,
                     navigation, logging, analytics adapters
```

Rule (identical to the backend's adapter pattern): **every external call —
HTTP, push, camera, Keychain/Keystore — sits behind a protocol/interface in
`domain/`, with the real implementation in `data/` and a fake/stub injected
in tests and previews.** Nothing in `features/*/view/` or `model/` imports
`URLSession`, `HttpInterface`, Firebase, or Keychain directly.

## Request flow

```
user gesture
  → Intent enum / callback into the feature model
  → model calls a use case → repository → data-source (API/cache)
  → model reduces the result into State (loading / success / failure)
  → View re-renders; errors surface as typed failures mapped to
    user-facing messages (backend `errorcodes` map to local error cases)
```

## Backend alignment

- The `status` / `service_id` / `message` / `data` envelope is decoded once,
  in the shared API client, not per-feature.
- `errorcodes` values map 1:1 to a typed error enum per platform — screens
  branch on the code, not on the message string.
- DTO names follow backend `models/` verbatim; both platforms generate or
  hand-maintain the same shape from the same source of truth.

## Previews, tests, quality gates

- Previews (Xcode) / Compose previews must build in CI without devices —
  they're the cheap equivalent of backend unit tests for UI state.
- Snapshot tests cover key screens only; the API layer is fake-injected so
  snapshots are deterministic.
- CI pipeline per app mirrors the backend shape: `build → unit-test →
  instrumented/e2e (opt-in) → archive → distribute (TestFlight / internal
  track)`.
