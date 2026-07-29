// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/renderer/html"
)

// deps holds data required for serving routes. Dependencies.
type deps struct {
	c *config.Config

	// staticDir should be path traversal attack resilient FS, such as the one
	// returned by "os.Root.FS".
	staticDir fs.FS

	templatesDir fs.FS

	// t is a compiled templates used by deps.Template().
	// Do not use this; use deps.Template() instead.
	t *template.Template

	markdown  html.MarkdownRenderer
	plaintext html.PlaintextRenderer
}

// template returns compiled HTML templates.
// Every routes should access templates through this method.
func (d *deps) template() *template.Template {
	if d.t != nil {
		return d.t
	}

	t := template.Must(template.ParseFS(d.templatesDir, "*"))
	return t
}

func (d *deps) write404(w http.ResponseWriter) {
	data := error404Data{
		Config: d.c,
	}

	w.WriteHeader(404)
	if err := d.template().ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) write500(w http.ResponseWriter) {
	data := error500Data{
		Config: d.c,
	}

	w.WriteHeader(500)
	if err := d.template().ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("500 template: %s", err)
	}
}
