# Project Intent

`ezgoblogthing` is an early-stage Go-first blog project. The repository has been initialized as a documentation-first workspace so the project's purpose, structure, and tools can be mapped before implementation choices are locked in.

## Working Purpose

Build a small, maintainable static blog generator. Authors should be able to put Markdown files in the repository and generate a complete static docroot that can be deployed to a host such as Netlify.

The first implementation should be written in Go and should avoid extra tooling unless the project clearly needs it.

Existing Hugo-style blog posts should be easy to move into this project. The goal is practical compatibility with common Markdown and front matter patterns, not a reimplementation of Hugo.

## Values

- Simplicity over machinery.
- Fast publishing over complex customization.
- Durable content formats over proprietary storage.
- Clear structure that future contributors and agents can navigate quickly.
- Effortless migration for straightforward Hugo posts.
- A boring deploy artifact: static files that can be served by common hosts.
- A small Go codebase before any frontend framework is introduced.

## Likely Use Cases

- Write posts as Markdown files.
- Bring existing Hugo Markdown posts into the project with minimal editing.
- Organize posts by date, topic, tag, or collection.
- Keep source content readable in the repository.
- Generate a public blog with predictable URLs.
- Generate posts under `/posts/`, pages under page-specific paths, and tag listing pages.
- Produce a deployable docroot for static hosting.

## Product Decisions

- The first version should be a static generator, not a dynamic server.
- Go is the primary implementation language.
- Plain Markdown is the expected content source format.
- The generated output should be a complete static docroot.
- The generated output directory is `dist/`.
- Blog posts should live under `/posts/` in the generated site.
- Pages and tags are part of the content model.
- Version 1 should skip files marked `draft: true`, but should not include a fuller drafts publishing workflow.
- Version 1 should not include feeds, search, image processing, or syntax highlighting.
- React may be used later, but it is not part of the initial scaffold.

## Open Product Decisions

- Exact content front matter schema beyond the Hugo-style fields seen in `examples/`.
- Exact URL model for pages and tags.
- Primary author/editor workflow.
- Deployment target, likely Netlify or another static host.

## Near-Term Success Criteria

The next useful milestone is a minimal working blog skeleton with:

- A Go module and one runnable generator command.
- A documented content model.
- A documented local serving command.
- A documented build command.
- At least one sample post proving the workflow.
- A generated docroot that can be served locally and deployed as static files.
- Example Hugo-style posts converted or copied into the real content tree with minimal edits.
