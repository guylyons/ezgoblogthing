package site

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Site struct {
	Posts []Item
	Pages []Item
	Tags  map[string][]Item
}

func Build(root string, out string) error {
	site, err := Load(root)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	if err := copyStatic(filepath.Join(root, "static"), out); err != nil {
		return err
	}
	if err := copyContentAssets(filepath.Join(root, "content", "posts"), filepath.Join(out, "posts")); err != nil {
		return err
	}
	if err := copyContentAssets(filepath.Join(root, "content", "content"), out); err != nil {
		return err
	}

	if err := writeIndex(out, site); err != nil {
		return err
	}
	for _, post := range site.Posts {
		if err := writeItem(out, filepath.Join("posts", post.Slug), post); err != nil {
			return err
		}
	}
	for _, page := range site.Pages {
		if err := writeItem(out, page.Slug, page); err != nil {
			return err
		}
	}
	for tag, posts := range site.Tags {
		if err := writeTag(out, tag, posts); err != nil {
			return err
		}
	}
	return nil
}

func Load(root string) (Site, error) {
	posts, err := loadItems(filepath.Join(root, "content", "posts"), KindPost)
	if err != nil {
		return Site{}, err
	}
	pages, err := loadItems(filepath.Join(root, "content", "pages"), KindPage)
	if err != nil {
		return Site{}, err
	}
	contentPages, err := loadContentPages(filepath.Join(root, "content", "content"))
	if err != nil {
		return Site{}, err
	}
	pages = append(pages, contentPages...)

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].SortDate.After(posts[j].SortDate)
	})
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Slug < pages[j].Slug
	})

	tags := map[string][]Item{}
	for _, item := range append(posts, pages...) {
		for _, tag := range item.Tags {
			tags[tag] = append(tags[tag], item)
		}
	}

	return Site{Posts: posts, Pages: pages, Tags: tags}, nil
}

func loadItems(dir string, kind Kind) ([]Item, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}
	var items []Item
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		if filepath.Base(path) == "_index.md" {
			return nil
		}
		item, err := ParseMarkdownFile(path, kind)
		if err != nil {
			return err
		}
		item.Slug = routeSlug(dir, path)
		if item.Draft {
			return nil
		}
		items = append(items, item)
		return nil
	})
	return items, err
}

func loadContentPages(dir string) ([]Item, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}
	var items []Item
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != dir && filepath.Base(path) == "posts" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		item, err := ParseMarkdownFile(path, KindPage)
		if err != nil {
			return err
		}
		item.Slug = routeSlug(dir, path)
		if item.Slug == "" || item.Draft {
			return nil
		}
		items = append(items, item)
		return nil
	})
	return items, err
}

func routeSlug(root string, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		rel = filepath.Base(path)
	}
	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	if filepath.Base(rel) == "_index" {
		rel = filepath.Dir(rel)
		if rel == "." {
			return ""
		}
	}
	return filepath.ToSlash(rel)
}

func writeIndex(out string, site Site) error {
	return writeTemplate(filepath.Join(out, "index.html"), indexTemplate, map[string]any{
		"Title": "ezgoblogthing",
		"Posts": site.Posts,
		"Pages": site.Pages,
	})
}

func writeItem(out string, route string, item Item) error {
	item.HTML = omitMissingLocalImages(out, item.HTML)
	return writeTemplate(filepath.Join(out, route, "index.html"), itemTemplate, item)
}

func writeTag(out string, tag string, posts []Item) error {
	return writeTemplate(filepath.Join(out, "tags", slugify(tag), "index.html"), tagTemplate, map[string]any{
		"Title": fmt.Sprintf("Posts tagged %s", tag),
		"Tag":   tag,
		"Posts": posts,
	})
}

func writeTemplate(path string, text string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmpl, err := template.New(filepath.Base(path)).Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
		"itemURL":  itemURL,
		"slug":     slugify,
		"date":     formatDate,
		"siteNav":  siteNav,
	}).Parse(text)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func itemURL(item Item) string {
	switch item.Kind {
	case KindPost:
		return "/posts/" + item.Slug + "/"
	default:
		return "/" + item.Slug + "/"
	}
}

var (
	localImageParagraphPattern = regexp.MustCompile(`<p><img src="(/[^"]+)" alt="[^"]*"></p>\n?`)
	localImagePattern          = regexp.MustCompile(`<img src="(/[^"]+)" alt="[^"]*">\n?`)
)

func omitMissingLocalImages(out string, html string) string {
	html = localImageParagraphPattern.ReplaceAllStringFunc(html, func(tag string) string {
		src := localImageParagraphPattern.FindStringSubmatch(tag)[1]
		if localAssetExists(out, src) {
			return tag
		}
		return ""
	})
	html = localImagePattern.ReplaceAllStringFunc(html, func(tag string) string {
		src := localImagePattern.FindStringSubmatch(tag)[1]
		if localAssetExists(out, src) {
			return tag
		}
		return ""
	})
	return html
}

func localAssetExists(out string, src string) bool {
	rel := strings.TrimPrefix(src, "/")
	if rel == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(out, filepath.FromSlash(rel)))
	return err == nil
}

func copyStatic(src string, out string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return copyFile(path, filepath.Join(out, rel))
	})
}

func copyContentAssets(src string, out string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != src && src == filepath.Join(filepath.Dir(src), "content") && filepath.Base(path) == "posts" {
				return filepath.SkipDir
			}
			return nil
		}
		switch {
		case filepath.Ext(path) == ".md":
			return nil
		case filepath.Ext(path) == ".html":
			return nil
		case strings.HasPrefix(filepath.Base(path), "."):
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return copyFile(path, filepath.Join(out, rel))
	})
}

func copyFile(src string, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var out strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			out.WriteRune(r)
			lastDash = false
		case !lastDash:
			out.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(out.String(), "-")
}

func siteNav() template.HTML {
	return template.HTML(`<nav class="top-nav" aria-label="Pages"><a href="/about/">About</a><span>*</span><a href="/links/">Links</a><span>*</span><a href="/contact/">Contact</a></nav>`)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	day := t.Day()
	suffix := "th"
	if day%100 < 11 || day%100 > 13 {
		switch day % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%s %d%s, %d", t.Month(), day, suffix, t.Year())
}

const indexTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="/style.css">
</head>
<body>
  <main>
    {{ siteNav }}
    <h1>{{ .Title }}</h1>
    <section>
      <h2>Posts</h2>
      {{ range .Posts }}
      <article>
        <h3><a href="/posts/{{ .Slug }}/">{{ .Title }}</a></h3>
        {{ with date .SortDate }}<p>{{ . }}</p>{{ end }}
      </article>
      {{ else }}
      <p>No posts yet.</p>
      {{ end }}
    </section>
    {{ if .Pages }}
    <section>
      <h2>Pages</h2>
      <ul class="page-list">{{ range .Pages }}<li><a href="/{{ .Slug }}/">{{ .Title }}</a></li>{{ end }}</ul>
    </section>
    {{ end }}
  </main>
</body>
</html>
`

const itemTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="/style.css">
</head>
<body>
  <main>
    {{ siteNav }}
    <p><a href="/">Home</a></p>
    <article>
      <h1>{{ .Title }}</h1>
      {{ with date .SortDate }}<p>{{ . }}</p>{{ end }}
      {{ if .Tags }}<p>{{ range .Tags }}<a href="/tags/{{ slug . }}/">{{ . }}</a> {{ end }}</p>{{ end }}
      {{ .HTML | safeHTML }}
    </article>
  </main>
</body>
</html>
`

const tagTemplate = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="/style.css">
</head>
<body>
  <main>
    {{ siteNav }}
    <p><a href="/">Home</a></p>
    <h1>{{ .Title }}</h1>
    {{ range .Posts }}
    <article>
      <h2><a href="{{ itemURL . }}">{{ .Title }}</a></h2>
      {{ with date .SortDate }}<p>{{ . }}</p>{{ end }}
    </article>
    {{ end }}
  </main>
</body>
</html>
`
