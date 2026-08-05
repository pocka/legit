// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package errors

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pocka/legit/core"
	"github.com/pocka/legit/templates"
)

// WriteNotFound writes 404 page response to "w".
func WriteNotFound(core *core.Core, w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)

	data := templates.Error404Data{
		Config:    core.Config,
		Diagnosis: fmt.Sprintf("Status: 404\nTimestamp: %s\nURL: %s\nUA: %s", now, r.URL, r.UserAgent()),
	}

	w.WriteHeader(http.StatusNotFound)
	if err := core.Template().ExecuteTemplate(w, "404", data); err != nil {
		log.Printf("execute template/404: %s", err)
	}
}

// WriteInternalServerError writes 500 page response to "w".
func WriteInternalServerError(core *core.Core, w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)

	data := templates.Error500Data{
		Config:    core.Config,
		Diagnosis: fmt.Sprintf("Status: 500\nTimestamp: %s\nURL: %s\nUA: %s", now, r.URL, r.UserAgent()),
	}

	w.WriteHeader(http.StatusInternalServerError)
	if err := core.Template().ExecuteTemplate(w, "500", data); err != nil {
		log.Printf("execute template/500: %s", err)
	}
}
