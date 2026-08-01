// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/git"
	"github.com/pocka/legit/renderer/html"
)

var (
	errRepositoryIsIgnored = errors.New("repository is ignored")
)

// deps holds data required for serving routes. Dependencies.
type deps struct {
	c *config.Config

	scanRoot *os.Root

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

// resolveRepository returns filepath to a repository.
func (d *deps) resolveRepository(name string) (string, error) {
	entry, err := d.scanRoot.Stat(name)
	if err != nil {
		return "", err
	}

	dirname := entry.Name()
	if d.isIgnored(dirname) {
		return "", errRepositoryIsIgnored
	}

	return filepath.Join(d.scanRoot.Name(), dirname), nil
}

func (d *deps) openRepository(name string, ref string) (repo *git.GitRepo, dirname string, err error) {
	path, err := d.resolveRepository(name)
	if err != nil {
		return nil, "", err
	}

	dirname = filepath.Base(path)
	repo, err = git.Open(path, ref)
	return
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

func (d *deps) write404(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Format(time.RFC3339)

	data := error404Data{
		Config:    d.c,
		Diagnosis: fmt.Sprintf("Status: 404\nTimestamp: %s\nURL: %s\nUA: %s", timestamp, r.URL, r.UserAgent()),
	}

	w.WriteHeader(404)
	if err := d.template().ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("404 template: %s", err)
	}
}

func (d *deps) write500(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Format(time.RFC3339)

	data := error500Data{
		Config:    d.c,
		Diagnosis: fmt.Sprintf("Status: 500\nTimestamp: %s\nURL: %s\nUA: %s", timestamp, r.URL, r.UserAgent()),
	}

	w.WriteHeader(500)
	if err := d.template().ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("500 template: %s", err)
	}
}
