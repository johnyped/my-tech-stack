# Web Frontend Tech Stack

## Stack

| Concern | Choice | Notes |
|---|---|---|
| Runtime + package manager | **bun** | `bun install` / `bunx` / `bun run` / `bun test` — one binary replaces node + npm + much of the jest/vitest ceremony |
| Bundler / dev server | **Vite** | `vite.config.ts`, `@vitejs/plugin-vue`; dev via `bun run dev`, production build via `bun run build` (`vue-tsc -b && vite build`) |
| Framework | **Vue 3** | SFCs with `<script setup lang="ts">`, Composition API only — Options API and `defineComponent` are legacy paths, not written by choice |
| Language | **TypeScript** | `strict: true`; `vue-tsc` in the build (no type-erasure-only "type checking" in editors) |
| Component library | **shadcn-vue** | shadcn/ui's component *source model* for Vue, on Reka UI headless primitives; components are vendored into `src/components/ui/` and owned by us, not installed as a black-box package |
| Styling | **Tailwind CSS v4** | via `@tailwindcss/vite`, CSS-first config (`@import "tailwindcss"`); theming through the shadcn CSS variable tokens, no runtime theme engines |
| Routing | **Vue Router** | history mode, route-level lazy loading (`() => import(...)`), typed `RouteMeta` for auth/permissions |
| State | **Pinia, opt-in** | composables + `provide/inject` first; Pinia only when state genuinely crosses unrelated views |
| HTTP | typed `src/lib/api/` wrapper | `fetch` (bun/Node-compatible) through one module: envelope decode, error mapping, auth headers |
| Test | `bun test` (unit) + Playwright (E2E) | same split as backend unit/integration: unit fast, E2E only for key flows |
| Lint/format | ESLint (`eslint-plugin-vue` + `@typescript-eslint`) + Prettier | single source: run via `bunx`; CI fails on lint errors |

## Why these and not the alternative

| Choice | Instead of | Why |
|---|---|---|
| Vue 3 | React, Svelte, Angular | team fluency + SFC ergonomics; the component story (shadcn-vue + Reka) and TypeScript support are mature enough to not need the React ecosystem |
| shadcn-vue | Element Plus / Naive UI / Vuetify as a *system* | shadcn's source-ownership model (you own the component code) keeps design tokens and Tailwind consistent and lets us eject a component without fighting a library's release cycle |
| Tailwind v4 | CSS Modules, styled-components | atomic utilities + the shadcn token variables; zero-runtime and no class-name collisions |
| bun | Node + npm/pnpm | speed of install/`bunx`/test and a single toolchain to document; Vite and Vue tooling are bun-compatible |
| Pinia opt-in | Vuex, always-on stores | most app state is server state + view state; a store layer is only added where it earns maintenance |

## Version posture

Pin to the newest stable of each at project start (bun latest, Vite latest,
Vue 3 latest, Tailwind v4, shadcn-vue latest, TS 5.x) and keep a
`DEPENDENCIES.md`-style note in the repo when intentionally staying behind.
No `--experimental` / preview features in app code.
