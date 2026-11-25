# Mobile Tech Stack

## iOS

| | |
|---|---|
| Language | Swift 6.x (strict concurrency / Swift 6 language mode) |
| UI | SwiftUI (SwiftUI App lifecycle — no Storyboards in new apps; UIKit interop via `UIViewControllerRepresentable` only when a component genuinely needs it) |
| Deployment target | Latest two iOS majors (SwiftUI's `@Observable` observation system requires iOS 17+ — that's the floor for new apps) |
| Tooling | Current stable Xcode; `xcodebuild`-driven CI (Xcode Cloud or GitLab CI macOS runners) |
| Dependencies | Swift Package Manager only — no CocoaPods for new apps |
| Testing | XCTest (unit) + view/flow tests via XCUITest or ViewKit-style tooling; snapshot tests (swift-snapshot-testing) for key screens |

## Android

| | |
|---|---|
| Language | Kotlin 2.x (K2 compiler, `kotlinx.coroutines` + Flow) |
| UI | Jetpack Compose (latest Compose BOM), Material 3 |
| Tooling | Latest stable Android Studio + AGP; Gradle Kotlin DSL with a version catalog (`gradle/libs.versions.toml`) |
| Build hygiene | `kotlin("plugin.compose")` + Compose compiler plugin via Gradle; KSP for codegen (Room, Hilt if/when added) |
| Persistence | Room for local cache; DataStore for key-value settings |
| Testing | JUnit + `androidx.compose.ui.test` (Compose test rule); Paparazzi/Roborazzi-style snapshots for key screens |

## Cross-platform choices (deliberate)

| | |
|---|---|
| No cross-platform UI toolkit | Flutter/React Native/Compose Multiplatform are not the default — the stack is explicitly native SwiftUI + native Compose |
| Networking | Native clients per platform: `URLSession` + `Codable` (iOS), Ktor Client (Android — shared-logic friendly) |
| Serialization | iOS: `Codable` (synthesised where possible). Android: `kotlinx.serialization` — JSON aligned field-for-field with the backend envelope |
| Auth | OAuth2 / bearer tokens in Keychain (iOS) / Keystore + EncryptedSharedPreferences (Android); never in app-level storage |
| API contracts | Source of truth is the backend's Bruno collection / OpenAPI spec — DTO field naming mirrors backend `models/` exactly so drift is visible |
| CI | Per-repo GitLab CI matching the backend pipeline shape (`build → test → archive/distribute`), fastlane only where an existing project already uses it |

## Why this and not the alternative

- **SwiftUI/Compose, not Flutter/RN**: the apps are form-heavy, API-driven
  product clients — platform look-and-feel and OS integration (push,
  notifications, widgets, biometrics) matter more than code sharing.
- **Native clients, not shared Kotlin/Swift logic layer**: at this team's
  app count, keeping one stack per platform is cheaper than maintaining a
  shared core with expect/actual or bridging seams.
- **`@Observable` + ViewModel, not VIPER/Clean-architecture ceremony**:
  keeps the same spirit as the backend's "thin handler, domain package owns
  the logic" rule — UI renders state, one observable model per screen owns
  it, data sources sit under that.
