package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Renderer struct {
	pages map[string]*template.Template
	base  *template.Template
	md    goldmark.Markdown
}

func New(templateFS fs.FS) (*Renderer, error) {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	r := &Renderer{md: md, pages: make(map[string]*template.Template)}

	funcMap := template.FuncMap{
		"markdown":      r.renderMarkdown,
		"formatDate":    formatDate,
		"today":         func() string { return time.Now().Format("2006-01-02") },
		"seq":           seq,
		"eq":            func(a, b any) bool { return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b) },
		"contains":      strings.Contains,
		"lower":         strings.ToLower,
		"upper":         strings.ToUpper,
		"replace":       func(s, old, new string) string { return strings.ReplaceAll(s, old, new) },
		"priorityLabel": priorityLabel,
		"statusLabel":   statusLabel,
		"categoryLabel": statusLabel,
		"add":           func(a, b int) int { return a + b },
		"sub":           func(a, b int) int { return a - b },
		"safeHTML":      func(s string) template.HTML { return template.HTML(s) },
	}

	// Find all page templates
	pageFiles, err := fs.Glob(templateFS, "templates/pages/*.html")
	if err != nil {
		return nil, fmt.Errorf("glob pages: %w", err)
	}

	for _, pagePath := range pageFiles {
		name := pagePath[len("templates/pages/"):]
		// Each page template = layout + all partials + this page
		t, err := template.New("").Funcs(funcMap).ParseFS(templateFS,
			"templates/layout.html",
			"templates/partials/*.html",
			pagePath,
		)
		if err != nil {
			return nil, fmt.Errorf("parse page %s: %w", name, err)
		}
		r.pages[name] = t
	}

	// Also parse partials alone for Partial() calls
	r.base, err = template.New("").Funcs(funcMap).ParseFS(templateFS,
		"templates/layout.html",
		"templates/partials/*.html",
		"templates/pages/*.html",
	)
	if err != nil {
		return nil, fmt.Errorf("parse base: %w", err)
	}

	return r, nil
}

func (r *Renderer) renderMarkdown(s string) template.HTML {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(s), &buf); err != nil {
		return template.HTML("<p>" + template.HTMLEscapeString(s) + "</p>")
	}
	return template.HTML(buf.String())
}

func (r *Renderer) RenderMarkdown(s string) string {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(s), &buf); err != nil {
		return s
	}
	return buf.String()
}

type PageData struct {
	Title   string
	Content any
	Flash   string
}

func (r *Renderer) Page(w http.ResponseWriter, status int, page string, data PageData) {
	t, ok := r.pages[page]
	if !ok {
		http.Error(w, "Template not found: "+page, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	io.Copy(w, &buf)
}

func (r *Renderer) Partial(w http.ResponseWriter, status int, name string, data any) {
	var buf bytes.Buffer
	err := r.base.ExecuteTemplate(&buf, name, data)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	io.Copy(w, &buf)
}

func formatDate(s string) string {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return s
	}
	return t.Format("02 Jan 2006")
}

func seq(start, end int) []int {
	var s []int
	for i := start; i <= end; i++ {
		s = append(s, i)
	}
	return s
}

func priorityLabel(p int) string {
	labels := map[int]string{1: "Critical", 2: "High", 3: "Medium", 4: "Low", 5: "Minimal"}
	if l, ok := labels[p]; ok {
		return l
	}
	return "Unknown"
}

func statusLabel(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	if len(s) == 0 {
		return s
	}
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
