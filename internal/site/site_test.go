package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMarkdownFileSupportsHugoStyleFrontMatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hello-world.md")
	err := os.WriteFile(path, []byte(`---
title: "Hello World"
date: 2019-01-06T07:01:11-05:00
draft: true
tags: ["go", poetry]
---

# Hello

This is a [link](https://example.com).
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	item, err := ParseMarkdownFile(path, KindPost)
	if err != nil {
		t.Fatal(err)
	}

	if item.Title != "Hello World" {
		t.Fatalf("Title = %q, want %q", item.Title, "Hello World")
	}
	if item.Slug != "hello-world" {
		t.Fatalf("Slug = %q, want %q", item.Slug, "hello-world")
	}
	if !item.Draft {
		t.Fatal("Draft = false, want true")
	}
	if strings.Join(item.Tags, ",") != "go,poetry" {
		t.Fatalf("Tags = %#v, want go, poetry", item.Tags)
	}
	if !strings.Contains(item.HTML, `<h1>Hello</h1>`) {
		t.Fatalf("HTML did not render heading: %s", item.HTML)
	}
	if !strings.Contains(item.HTML, `<a href="https://example.com">link</a>`) {
		t.Fatalf("HTML did not render link: %s", item.HTML)
	}
}

func TestParseMarkdownFileSupportsTOMLFrontMatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "design-4-drupal-boston.md")
	err := os.WriteFile(path, []byte(`+++
date = "2019-06-27T04:00:00+00:00"
tags = ["drupal"]
title = "Design 4 Drupal Boston"
+++

Body text.
`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	item, err := ParseMarkdownFile(path, KindPost)
	if err != nil {
		t.Fatal(err)
	}

	if item.Title != "Design 4 Drupal Boston" {
		t.Fatalf("Title = %q", item.Title)
	}
	if strings.Join(item.Tags, ",") != "drupal" {
		t.Fatalf("Tags = %#v, want drupal", item.Tags)
	}
}

func TestBuildGeneratesDistWithPostsPagesTagsAndStaticAssets(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "posts", "hello.md"), `---
title: Hello
date: 2026-06-16
tags: ["go"]
---
Hello post.`)
	mustWrite(t, filepath.Join(root, "content", "pages", "about.md"), `---
title: About
---
About page.`)
	mustWrite(t, filepath.Join(root, "content", "pages", "links.md"), `---
title: Links
---
Links page.`)
	mustWrite(t, filepath.Join(root, "content", "pages", "contact.md"), `---
title: Contact
---
Contact page.`)
	mustWrite(t, filepath.Join(root, "static", "style.css"), `body { color: #222; }`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, filepath.Join(out, "index.html"), "Hello")
	assertFileContains(t, filepath.Join(out, "index.html"), `<nav class="top-nav" aria-label="Pages"><a href="/about/">About</a><span>*</span><a href="/links/">Links</a><span>*</span><a href="/contact/">Contact</a></nav>`)
	assertFileContains(t, filepath.Join(out, "posts", "hello", "index.html"), "Hello post.")
	assertFileContains(t, filepath.Join(out, "posts", "hello", "index.html"), "June 16th, 2026")
	assertFileContains(t, filepath.Join(out, "about", "index.html"), "About page.")
	assertFileContains(t, filepath.Join(out, "links", "index.html"), "Links page.")
	assertFileContains(t, filepath.Join(out, "contact", "index.html"), "Contact page.")
	assertFileContains(t, filepath.Join(out, "tags", "go", "index.html"), "Hello")
	assertFileContains(t, filepath.Join(out, "style.css"), "color: #222")
}

func TestBuildPreservesNestedPostRoutesAndCopiesPostAssets(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "posts", "2024", "keyboardio.md"), `---
title: Keyboardio
date: 2024-01-02
---
![keyboard](/posts/2024/keyboardio.png)`)
	mustWriteBytes(t, filepath.Join(root, "content", "posts", "2024", "keyboardio.png"), []byte("image bytes"))

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, filepath.Join(out, "index.html"), `/posts/2024/keyboardio/`)
	assertFileContains(t, filepath.Join(out, "posts", "2024", "keyboardio", "index.html"), `<img src="/posts/2024/keyboardio.png" alt="keyboard">`)
	assertFileContains(t, filepath.Join(out, "posts", "2024", "keyboardio.png"), "image bytes")
}

func TestBuildOmitsMissingLocalImages(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "posts", "missing-image.md"), `---
title: Missing Image
---
![missing](/uploads/missing.jpg)

![remote](https://example.com/image.jpg)

Body.`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(filepath.Join(out, "posts", "missing-image", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if strings.Contains(got, `/uploads/missing.jpg`) {
		t.Fatalf("missing local image was rendered:\n%s", got)
	}
	if !strings.Contains(got, `https://example.com/image.jpg`) {
		t.Fatalf("remote image was not preserved:\n%s", got)
	}
}

func TestBuildGeneratesSectionPagesFromHugoContentDirectory(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "content", "recipes", "_index.md"), `---
title: Recipes
---
Recipe collection.`)
	mustWrite(t, filepath.Join(root, "content", "content", "recipes", "bagels.md"), `---
title: Bagels
date: 2026-02-24
tags: ["recipes"]
---
Bagel recipe.`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	assertFileContains(t, filepath.Join(out, "index.html"), `/recipes/`)
	assertFileContains(t, filepath.Join(out, "recipes", "index.html"), "Recipe collection.")
	assertFileContains(t, filepath.Join(out, "recipes", "bagels", "index.html"), "Bagel recipe.")
	assertFileContains(t, filepath.Join(out, "tags", "recipes", "index.html"), "Bagels")
	assertFileContains(t, filepath.Join(out, "tags", "recipes", "index.html"), `/recipes/bagels/`)
}

func TestBuildSkipsDraftContent(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "posts", "draft.md"), `---
title: Draft
draft: true
---
Draft body.`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(out, "posts", "draft", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("draft post was generated, stat err: %v", err)
	}
}

func TestRenderMarkdownKeepsSimpleRawHTMLAndRendersLists(t *testing.T) {
	got := renderMarkdown(`<img src="/posts/2024/keyboardio.png" alt="Keyboardio">

## Resources

- [Keyboardio](https://www.keyboard.io)
- Plain item`)

	if !strings.Contains(got, `<img src="/posts/2024/keyboardio.png" alt="Keyboardio">`) {
		t.Fatalf("raw HTML was not preserved:\n%s", got)
	}
	if !strings.Contains(got, `<ul>`) || !strings.Contains(got, `<li><a href="https://www.keyboard.io">Keyboardio</a></li>`) {
		t.Fatalf("list was not rendered:\n%s", got)
	}
}

func TestRenderMarkdownHandlesAdjacentHeadingsAndOrderedLists(t *testing.T) {
	got := renderMarkdown(`### Optional toppings
- Sesame
- Poppy

## Instructions
1. **Mix:** Combine ingredients.
2. **Bake:** Bake until done.`)

	if !strings.Contains(got, "<h3>Optional toppings</h3>") {
		t.Fatalf("heading was not isolated:\n%s", got)
	}
	if !strings.Contains(got, "<li>Sesame</li>") || !strings.Contains(got, "<ol>") {
		t.Fatalf("lists were not rendered:\n%s", got)
	}
	if !strings.Contains(got, "<li><strong>Mix:</strong> Combine ingredients.</li>") {
		t.Fatalf("ordered list inline markup was not rendered:\n%s", got)
	}
}

func TestRenderMarkdownSupportsYouTubeShortcode(t *testing.T) {
	got := renderMarkdown(`{{< youtube 4jc1TJoNUuA >}}`)

	if !strings.Contains(got, `https://www.youtube.com/embed/4jc1TJoNUuA`) {
		t.Fatalf("youtube shortcode was not rendered:\n%s", got)
	}
	if !strings.Contains(got, `<iframe`) {
		t.Fatalf("youtube embed iframe was not rendered:\n%s", got)
	}
}

func TestBuildDoesNotCopyHiddenContentJunk(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "posts", "hello.md"), `---
title: Hello
---
Hello.`)
	mustWrite(t, filepath.Join(root, "content", "posts", ".DS_Store"), `junk`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(out, "posts", ".DS_Store")); !os.IsNotExist(err) {
		t.Fatalf("hidden content junk was copied, stat err: %v", err)
	}
}

func TestBuildDoesNotCopyUnsupportedContentHTML(t *testing.T) {
	root := t.TempDir()
	mustWrite(t, filepath.Join(root, "content", "content", "_index.html"), `<img src="/missing.png">`)

	out := filepath.Join(root, "dist")
	if err := Build(root, out); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(out, "_index.html")); !os.IsNotExist(err) {
		t.Fatalf("unsupported content HTML was copied, stat err: %v", err)
	}
}

func mustWrite(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustWriteBytes(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertFileContains(t *testing.T, path string, want string) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), want) {
		t.Fatalf("%s did not contain %q:\n%s", path, want, string(b))
	}
}
