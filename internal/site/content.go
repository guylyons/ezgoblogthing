package site

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Kind string

const (
	KindPost Kind = "post"
	KindPage Kind = "page"
)

type Item struct {
	Kind     Kind
	Title    string
	Date     string
	SortDate time.Time
	Slug     string
	Tags     []string
	Draft    bool
	HTML     string
}

func ParseMarkdownFile(path string, kind Kind) (Item, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Item{}, err
	}

	meta, body, err := splitFrontMatter(string(b))
	if err != nil {
		return Item{}, fmt.Errorf("%s: %w", path, err)
	}

	item := Item{
		Kind:  kind,
		Slug:  strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Title: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
	}
	applyFrontMatter(&item, meta)
	item.SortDate = parseDate(item.Date)
	item.HTML = renderMarkdown(body)

	return item, nil
}

func splitFrontMatter(input string) (map[string]string, string, error) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	if strings.HasPrefix(input, "---\n") {
		return parseDelimitedFrontMatter(input, "---", parseYAMLish)
	}
	if strings.HasPrefix(input, "+++\n") {
		return parseDelimitedFrontMatter(input, "+++", parseTOMLish)
	}
	return map[string]string{}, input, nil
}

func parseDelimitedFrontMatter(input string, delimiter string, parser func(string) map[string]string) (map[string]string, string, error) {
	lines := strings.Split(input, "\n")
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == delimiter {
			meta := parser(strings.Join(lines[1:i], "\n"))
			body := strings.Join(lines[i+1:], "\n")
			return meta, body, nil
		}
	}
	return nil, "", fmt.Errorf("unterminated %s front matter", delimiter)
}

func parseYAMLish(input string) map[string]string {
	meta := map[string]string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		meta[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return meta
}

func parseTOMLish(input string) map[string]string {
	meta := map[string]string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		meta[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return meta
}

func applyFrontMatter(item *Item, meta map[string]string) {
	if title := cleanScalar(meta["title"]); title != "" {
		item.Title = title
	}
	if date := cleanScalar(meta["date"]); date != "" {
		item.Date = date
	}
	if draft := strings.ToLower(cleanScalar(meta["draft"])); draft == "true" {
		item.Draft = true
	}
	item.Tags = parseTags(meta["tags"])
}

func cleanScalar(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, `"`)
	value = strings.Trim(value, `'`)
	return value
}

func parseTags(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	value = strings.TrimPrefix(value, "[")
	value = strings.TrimSuffix(value, "]")
	parts := strings.Split(value, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := cleanScalar(strings.TrimSpace(part))
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	sort.Strings(tags)
	return tags
}

func parseDate(value string) time.Time {
	value = cleanScalar(value)
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000Z07:00",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

var (
	imagePattern       = regexp.MustCompile(`!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	linkPattern        = regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	youtubePattern     = regexp.MustCompile(`^\{\{<\s*youtube\s+([A-Za-z0-9_-]+)\s*>\}\}$`)
	orderedItemPattern = regexp.MustCompile(`^\d+\.\s+(.+)$`)
	strongPattern      = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	emPattern          = regexp.MustCompile(`\*([^*]+)\*`)
)

func renderMarkdown(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	blocks := strings.Split(input, "\n\n")
	var out bytes.Buffer
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		if youtube := renderYouTubeShortcode(block); youtube != "" {
			out.WriteString(youtube)
			continue
		}
		if isRawHTMLBlock(block) {
			out.WriteString(block)
			out.WriteString("\n")
			continue
		}
		renderBlock(&out, block)
	}
	return out.String()
}

func renderBlock(out *bytes.Buffer, block string) {
	var paragraph []string
	listKind := ""

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(renderInline(strings.Join(paragraph, " ")))
		out.WriteString("</p>\n")
		paragraph = nil
	}
	flushList := func() {
		if listKind == "" {
			return
		}
		out.WriteString("</")
		out.WriteString(listKind)
		out.WriteString(">\n")
		listKind = ""
	}
	startList := func(kind string) {
		flushParagraph()
		if listKind == kind {
			return
		}
		flushList()
		out.WriteString("<")
		out.WriteString(kind)
		out.WriteString(">\n")
		listKind = kind
	}

	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			flushParagraph()
			flushList()
			continue
		}
		switch {
		case strings.HasPrefix(line, "### "):
			flushParagraph()
			flushList()
			writeHeading(out, "h3", strings.TrimSpace(strings.TrimPrefix(line, "### ")))
		case strings.HasPrefix(line, "## "):
			flushParagraph()
			flushList()
			writeHeading(out, "h2", strings.TrimSpace(strings.TrimPrefix(line, "## ")))
		case strings.HasPrefix(line, "# "):
			flushParagraph()
			flushList()
			writeHeading(out, "h1", strings.TrimSpace(strings.TrimPrefix(line, "# ")))
		case strings.HasPrefix(line, "- "):
			startList("ul")
			writeListItem(out, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		case orderedItemPattern.MatchString(line):
			startList("ol")
			matches := orderedItemPattern.FindStringSubmatch(line)
			writeListItem(out, matches[1])
		default:
			flushList()
			paragraph = append(paragraph, line)
		}
	}
	flushParagraph()
	flushList()
}

func writeHeading(out *bytes.Buffer, tag string, text string) {
	out.WriteString("<")
	out.WriteString(tag)
	out.WriteString(">")
	out.WriteString(renderInline(text))
	out.WriteString("</")
	out.WriteString(tag)
	out.WriteString(">\n")
}

func writeListItem(out *bytes.Buffer, text string) {
	out.WriteString("<li>")
	out.WriteString(renderInline(text))
	out.WriteString("</li>\n")
}

func renderInline(input string) string {
	escaped := html.EscapeString(input)
	escaped = imagePattern.ReplaceAllString(escaped, `<img src="$2" alt="$1">`)
	escaped = linkPattern.ReplaceAllString(escaped, `<a href="$2">$1</a>`)
	escaped = strongPattern.ReplaceAllString(escaped, `<strong>$1</strong>`)
	escaped = emPattern.ReplaceAllString(escaped, `<em>$1</em>`)
	return escaped
}

func renderYouTubeShortcode(block string) string {
	matches := youtubePattern.FindStringSubmatch(strings.TrimSpace(block))
	if matches == nil {
		return ""
	}
	id := matches[1]
	return fmt.Sprintf(`<div class="embed-container"><iframe src="https://www.youtube.com/embed/%s" title="YouTube video player" loading="lazy" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></div>
`, html.EscapeString(id))
}

func isRawHTMLBlock(block string) bool {
	block = strings.TrimSpace(block)
	return strings.HasPrefix(block, "<") && strings.HasSuffix(block, ">")
}
