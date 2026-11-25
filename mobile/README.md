# Mobile Blueprint

Native mobile apps for iOS and Android. One converged standard for both
platforms: native UI toolkits (SwiftUI / Jetpack Compose), declarative-first,
MVVM-shaped layering, and the same API contracts the backend blueprint
exposes (`status` / `service_id` / `message` / `data` envelope).

## Start here

| I want to... | Go to |
|---|---|
| See the language / framework / library choices and why | [document/TECH_STACK.md](document/TECH_STACK.md) |
| Understand the layering, data flow, and shared conventions across iOS/Android | [document/ARCHITECTURE.md](document/ARCHITECTURE.md) |

## New app bootstrap (cheat sheet)

```text
iOS:     Xcode → new app (SwiftUI App lifecycle) → Swift 6 language mode →
         SPM only; no CocoaPods in new projects.
Android: Android Studio new project → "Single Module App (Compose)" →
         Kotlin DSL + version catalogs (libs.versions.toml).
```

Both platforms start with: API client + auth wiring, feature-first module
layout, and previews/snapshots enabled from the first commit.
