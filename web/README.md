# Web Frontend Blueprint

The web stack: **Vite + bun + Vue 3 + TypeScript + shadcn UI + Tailwind CSS**.
Written in the same blueprint style as `backend/` — every choice here has a
reason, and new apps start from the bootstrap below rather than from scratch.

## Start here

| I want to... | Go to |
|---|---|
| See the stack choices and why | [document/TECH_STACK.md](document/TECH_STACK.md) |
| Understand project layout and data flow | [document/ARCHITECTURE.md](document/ARCHITECTURE.md) |

## New app bootstrap

```bash
bun create vite my-app --template vue-ts
cd my-app && bun install

# Tailwind CSS v4 (CSS-first, via the Vite plugin)
bun add tailwindcss @tailwindcss/vite
#   vite.config.ts → plugins: [vue(), tailwindcss()]
#   src/assets/main.css → @import "tailwindcss";
#   tsconfig: paths { "@/*": ["./src/*"] }

# shadcn-vue (shadcn/ui components for Vue, on Reka UI primitives)
bunx shadcn-vue@latest init

# Routing (default) + state (only if genuinely needed)
bun add vue-router
# bun add pinia   # skip until a real cross-view state need shows up
```

Conventions carried over from the backend blueprint: typed API layer against
the `status` / `service_id` / `message` / `data` envelope, `bun test` for
unit tests, Playwright for the small set of E2E flows.
