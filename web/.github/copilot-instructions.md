<!-- Helpful instructions for AI coding agents working in this repo -->
# Copilot / AI Agent Quick Instructions

Purpose: concise, actionable facts that let AI agents be productive immediately in this monorepo.

- Repo type: a PNPM monorepo (see `package.json` root and `pnpm-workspace.yaml`). Use `pnpm` (>=8.10.0) and Node >=18.12.0.
- App surface: primary frontend in `src/` (Vue 3 + Vite) and small auxiliary apps under `apps/` (notably `apps/test-server` for local mock APIs).

Quick dev commands
- Install: `pnpm install` (alias: `pnpm run bootstrap`). Root `preinstall` enforces pnpm only.
- Dev frontend: `pnpm dev` (runs `vite`). This opens a browser and uses the proxies defined in `vite.config.ts`.
- Build frontend: `pnpm run build:production` (or `pnpm run build:docker` / `build:test`). `postbuild` runs `node ./build/build.js`.
- Preview: `pnpm run preview` (build + `vite preview`).
- Start mock/test server: `cd apps/test-server && pnpm run start` (nodemon). For production use pm2: `pnpm run prod` in that folder.

Key files and where to look (examples)
- Vite + proxy config: [vite.config.ts](vite.config.ts#L1-L200) — dev server proxies many `/api/v1/*` routes to backend services on localhost ports 8088–8096. Update here for API surface changes.
- Frontend bootstrap: [src/main.ts](src/main.ts#L1-L120) — app setup, global plugins, and the `ChatBotPlugin` example; defaultReply shown here.
- Root scripts, dependencies, pnpm constraints: [package.json](package.json#L1-L60)
- Local test server: [apps/test-server/package.json](apps/test-server/package.json#L1-L60) and [apps/test-server/README.md](apps/test-server/README.md#L1-L12)

Monorepo conventions and patterns
- Local packages referenced with `workspace:*` (e.g., `@vben/hooks`, `@vben/vite-config`). Inspect `packages/` and `internal/` for shared configs and plugins.
- Vite config is extended via `@vben/vite-config` — prefer editing project-level `vite.config.ts` for app-specific overrides rather than altering the shared package.
- Proxy pattern: frontend relies on dev-time proxies to route to many backend services. Tests or local development expect those backend ports to be available or to use the provided `apps/test-server` mocks.
- Post-build step: `postbuild` executes `build/build.js` — CI and packaging rely on this artifact.

Code and style conventions
- TypeScript + Vue 3 composition API. Types are checked with `vue-tsc` via `type:check` script.
- Linting: `turbo run lint` orchestrates project linters. Individual tasks: `lint:eslint`, `lint:stylelint`, `lint:prettier` in root `package.json`.
- Commits: `husky` + `commitizen` via `cz-git` (see `prepare` and `config.commitizen`).

Testing and debugging notes for agents
- There is no centralized test harness in the root; use component/unit tests where present and run `pnpm` scripts per-package.
- To reproduce front-end behavior locally: run `pnpm dev` then inspect proxy target availability. If backend unavailable, run `apps/test-server` to provide basic endpoints (uploads, websockets, login mocks).

Integration points and external dependencies
- Multiple backend services expected by the frontend—see proxy targets in `vite.config.ts` (ports 8088–8096 and others). Modify these for local debugging.
- The chat assistant is wired as a plugin (`src/chat-bot`) and used in `src/main.ts`. Look under `src/chat-bot` for implementation details and examples.

What agents should change vs avoid
- Safe to change: feature code under `src/`, components, route logic, and local `apps/*` mock servers.
- Prefer NOT to change: shared config packages in `packages/` and `internal/` unless necessary — they are consumed by other workspaces.
- When in doubt about build or packaging, inspect `postbuild` (`build/build.js`) and the monorepo `turbo` pipeline before refactoring build steps.

Example actionable tasks (with file pointers)
- Update a proxy or add a mock route: edit [vite.config.ts](vite.config.ts#L1-L120) and add mocks under `mock/` or `apps/test-server`.
- Add global UI behavior or plugin: update [src/main.ts](src/main.ts#L1-L120) and component registration in `components/registerGlobComp.ts`.
- Modify shared TypeScript config: edit `internal/ts-config` or `packages/*` and run `pnpm install && pnpm -w -F <pkg> build` if you change package outputs.

If anything here is unclear or you want more detail (examples, common TODOs, or CI notes), tell me which area to expand and I will iterate.
