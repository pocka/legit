// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package pages

import (
	"net/http"
	"path"
	"strings"

	"github.com/pocka/legit/core"
	"github.com/pocka/legit/pages/debug"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/pages/repo"
)

const (
	staticPrefix = "/static/"
)

type Pages struct {
	core *core.Core
}

func New(core *core.Core) *Pages {
	return &Pages{core: core}
}

func (pages *Pages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if debug.Handle(w, r) {
		return
	}

	if r.URL.Path == "/" {
		pages.index(w, r)
		return
	}

	// "/static/<path>*"
	if strings.Index(r.URL.Path, staticPrefix) == 0 {
		http.ServeFileFS(w, r, pages.core.StaticDir, strings.TrimPrefix(r.URL.Path, staticPrefix))
		return
	}

	// Everything else will be treated as a repository page.
	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(segments) == 0 {
		errors.WriteNotFound(pages.core, w, r)
		return
	}

	repoName := segments[0]
	repoPath, err := pages.core.RepositoryPath(repoName)
	if err != nil {
		errors.WriteNotFound(pages.core, w, r)
		return
	}

	innerPath := path.Join(segments[1:]...)

	repo, err := repo.New(repoPath, "/"+innerPath, pages.core)
	if err != nil {
		errors.WriteNotFound(pages.core, w, r)
		return
	}

	repo.ServeHTTP(w, r)
}
