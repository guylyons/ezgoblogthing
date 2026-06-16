# ezgoblogthing

`ezgoblogthing` is a planning-stage Go static blog generator.

The intended workflow is:

1. Write blog content as Markdown in the repository.
2. Run a small Go generator.
3. Produce `dist/`, a complete static docroot for deployment to a host such as Netlify.

The project should make straightforward Hugo-style Markdown posts easy to bring over, while staying much smaller than Hugo itself.

No runnable implementation exists yet. Current project decisions live in:

- [docs/intent.md](docs/intent.md)
- [docs/content.md](docs/content.md)
- [docs/structure.md](docs/structure.md)
- [docs/tools.md](docs/tools.md)
