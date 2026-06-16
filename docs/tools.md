# Tools

Go is the selected primary implementation tool. This file should become the source of truth for local commands, dependencies, verification, and deployment steps as the generator is built.

## Current Tooling State

- Primary language/runtime: Go.
- Package manager: Go modules once `go.mod` exists.
- Framework or static site generator: custom minimal Go generator.
- Test runner: Go's built-in `go test` once implementation exists.
- Formatter: Go's built-in `gofmt` once implementation exists.
- Frontend tooling: not selected; React is a possible future addition.
- Deployment target: not selected; static hosting such as Netlify is likely.

## Tooling Principles

- Choose the smallest stack that supports the desired publishing workflow.
- Prefer Go standard-library capabilities before adding dependencies.
- Add dependencies only when they remove meaningful project complexity.
- Document every command a contributor needs to run locally.
- Keep generated output reproducible from source files.

## Tooling Direction

The current direction is:

- A Go command reads Markdown from `content/`.
- Templates and static assets are kept in the repository.
- The generator writes a complete static docroot to `dist/`.
- The generated docroot can be served locally or uploaded to a static host.
- Existing Hugo-style Markdown posts should be accepted when they use the supported front matter subset.
- React remains a future option for interactive behavior, not an initial dependency.

## Commands

No project commands exist yet because the Go module has not been scaffolded.

Expected future commands should follow this shape:

```sh
go run ./cmd/ezgoblogthing build
```

Purpose:

- Generate the static docroot from repository content.
- Run when previewing or preparing a deploy.
- Expected output: a generated directory such as `dist/`.
- Version 1 output should include `/posts/`, generated pages, and tag pages.

```sh
go test ./...
```

Purpose:

- Run all Go tests.
- Run before committing implementation changes.

```sh
gofmt -w ./cmd ./internal
```

Purpose:

- Format Go source files once those directories exist.
- Run after Go code changes.

## Verification

No automated verification exists yet because no implementation exists.

Once implementation begins, add the expected checks here:

- Formatting with `gofmt`.
- Tests with `go test ./...`.
- Static generation command.
- Link or content validation if the generator adds it.
- Import checks against representative files in `examples/`.

## Deployment

No deployment target has been selected yet.

Document the deployment path once chosen, including:

- Hosting provider.
- Build command.
- Output directory.
- Required environment variables.
- Preview and production publishing workflow.

Current deployment assumption:

- The output should be a static docroot that can be served by Netlify or a similar static host.
- The initial project should not depend on a hosting provider SDK or service-specific build tool.

## Version Control

This repository is initialized with Jujutsu (`jj`) using a colocated Git repository. Prefer `jj` for local version control workflows; use Git directly only when interoperability requires it.
