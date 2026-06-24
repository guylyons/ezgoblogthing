# Content Model

`ezgoblogthing` should make simple Markdown publishing feel direct: write files, keep them readable, and generate predictable static URLs.

## Source Content

The expected source directories are:

- `content/posts/`: blog posts.
- `content/pages/`: standalone pages.
- `content/content/`: migrated Hugo-style section content. Markdown outside `content/content/posts/` is treated as page content.
- `static/`: assets copied into the generated docroot.
- `examples/`: reference Hugo-style posts and assets used to test migration behavior.

Files in `examples/` are not production content by default.

## Front Matter

Version 1 should support the common front matter patterns visible in `examples/`:

- YAML front matter delimited by `---`.
- TOML front matter delimited by `+++`.
- `title`: display title.
- `date`: publish date, including date-only values and timestamp values with offsets.
- `draft`: accepted for compatibility. Files with `draft: true` are skipped during generation.
- `tags`: list of tag names.

The generator should tolerate straightforward Hugo-style Markdown posts when they use this supported subset. It should not try to implement Hugo's full content model, full shortcode system, taxonomy behavior, template lookup, or configuration surface.

Version 1 does support one narrow Hugo-style shortcode that appears in the sample content: `{{< youtube VIDEO_ID >}}`, which renders to a plain YouTube `<iframe>` embed.

## Routing

Version 1 should generate:

- Blog posts under `/posts/`, preserving nested source paths such as `content/posts/2024/example.md` to `/posts/2024/example/`.
- Pages under page-specific paths. `content/pages/about.md` becomes `/about/`; `content/content/recipes/bagels.md` becomes `/recipes/bagels/`.
- Section index pages from `_index.md`, such as `content/content/recipes/_index.md` to `/recipes/`.
- Tag listing pages.
- Non-Markdown content assets copied beside their generated routes, such as `content/posts/2024/image.png` to `/posts/2024/image.png`.

Open details:

- Whether page URLs are controlled by filename, front matter, or both.
- Tag URL shape and tag normalization rules. Version 1 uses `/tags/<normalized-tag>/`.

## Version 1 Non-Goals

Version 1 should not include:

- RSS or Atom feeds.
- Site search.
- Image processing.
- Syntax highlighting.
- Scheduled publishing.
- A full drafts publishing workflow beyond skipping `draft: true`.
