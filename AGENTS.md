# AGENTS.md

Guidance for AI coding agents working in this project.

## Project Intent

This repository is for `ezgoblogthing`, a Go-first static blog generator whose implementation has not been scaffolded yet. Treat this repo as a planning-first workspace until the first Go module and runnable command are added.

Primary goals:

- Keep the project easy to understand and modify.
- Document intent before adding framework, frontend, or build complexity.
- Prefer small, explicit files over broad abstractions.
- Preserve decisions in `docs/` as the project takes shape.
- Use Markdown as the durable source format for blog content.
- Make existing Hugo-style blog posts easy to move into the project, without trying to replicate Hugo.
- Generate a complete static docroot that can be shipped to a host such as Netlify.

See [docs/intent.md](docs/intent.md) for the current product direction.

## Current Structure

The repository currently contains planning documentation and baseline hygiene files:

```text
.
├── .gitignore
├── AGENTS.md
├── README.md
├── examples/
│   └── 2019/
└── docs/
    ├── content.md
    ├── intent.md
    ├── structure.md
    └── tools.md
```

See [docs/structure.md](docs/structure.md) for the intended layout as implementation begins.

## Agent Workflow

Before making changes:

1. Read this file and the relevant file in `docs/`.
2. Inspect the current tree; do not assume the repo still matches this document.
3. Keep edits scoped to the user's request.
4. Avoid introducing a framework, dependency, service, or build step without a clear reason.
5. Update docs when a decision changes project intent, structure, or tooling.

When adding implementation:

- Follow the structure documented in `docs/structure.md`, or update that document first if the structure needs to change.
- Keep generated assets, build outputs, and dependency folders out of source control.
- Add verification steps to `docs/tools.md` once a test, lint, or build command exists.
- Prefer Go standard-library solutions first. Add dependencies only when they remove real complexity.
- Do not add React, Node tooling, or frontend build steps until there is a concrete interactive feature that needs them.
- Preserve compatibility with common Hugo Markdown/front matter patterns when practical, but do not add broad Hugo behavior unless the project needs it.
- See [docs/content.md](docs/content.md) before changing Markdown parsing, routing, or content fields.

## Tooling Baseline

Go is the selected primary implementation language. No frontend package manager, JavaScript framework, test runner beyond Go's built-in tooling, or deployment target has been selected yet.

See [docs/tools.md](docs/tools.md) for the current tooling map and open decisions.

## Version Control

This repository is initialized with Jujutsu (`jj`) using a colocated Git repository. Prefer `jj` for local version control work; use Git directly only when interoperability requires it.

## Documentation Expectations

Use plain Markdown for project docs. Prefer short, decision-oriented sections:

- What is true now.
- Why it is true.
- What should happen next.

Avoid stale roadmap promises. If a section becomes speculative, label it as an open decision.
