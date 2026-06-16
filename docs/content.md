# Content Model

`ezgoblogthing` should make simple Markdown publishing feel direct: write files, keep them readable, and generate predictable static URLs.

## Source Content

The expected source directories are:

- `content/posts/`: blog posts.
- `content/pages/`: standalone pages.
- `static/`: assets copied into the generated docroot.
- `examples/`: reference Hugo-style posts and assets used to test migration behavior.

Files in `examples/` are not production content by default.

## Front Matter

Version 1 should support the common front matter patterns visible in `examples/`:

- YAML front matter delimited by `---`.
- TOML front matter delimited by `+++`.
- `title`: display title.
- `date`: publish date, including date-only values and timestamp values with offsets.
- `draft`: accepted for compatibility. Version 1 does not need a full drafts workflow, but generation behavior for `draft: true` still needs a concrete decision.
- `tags`: list of tag names.

The generator should tolerate straightforward Hugo-style Markdown posts when they use this supported subset. It should not try to implement Hugo's full content model, shortcode system, taxonomy behavior, template lookup, or configuration surface.

## Routing

Version 1 should generate:

- Blog posts under `/posts/`.
- Pages under page-specific paths.
- Tag listing pages.

Open details:

- Whether post URLs include year/month segments or only slugs.
- Whether page URLs are controlled by filename, front matter, or both.
- Tag URL shape and tag normalization rules.

## Version 1 Non-Goals

Version 1 should not include:

- RSS or Atom feeds.
- Site search.
- Image processing.
- Syntax highlighting.
- Scheduled publishing.
- A full drafts publishing workflow.
