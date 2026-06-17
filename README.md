# ezgoblogthing

`ezgoblogthing` is a small Go static blog generator.

![Generated ezgoblogthing site preview](github.png)

The intended workflow is:

1. Write blog content as Markdown in the repository.
2. Run a small Go generator.
3. Produce `dist/`, a complete static docroot for deployment to a host such as Netlify.

The project should make straightforward Hugo-style Markdown posts easy to bring over, while staying much smaller than Hugo itself.

## Commands

Build the static site:

```sh
go run ./cmd/ezgoblogthing build
```

Build and serve the static site locally:

```sh
go run ./cmd/ezgoblogthing serve
```

Then open <http://localhost:8080>.

Run tests:

```sh
go test ./...
```

Current project decisions live in:

- [docs/intent.md](docs/intent.md)
- [docs/content.md](docs/content.md)
- [docs/structure.md](docs/structure.md)
- [docs/tools.md](docs/tools.md)
