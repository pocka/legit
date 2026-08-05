// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"compress/gzip"
	"log"
	"net/http"

	"github.com/pocka/legit/git/service"
)

func (repo *Repo) uploadpack(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/x-git-upload-pack-result")
	header.Set("Connection", "Keep-Alive")
	header.Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)

	var err error
	reader := r.Body
	if r.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Printf("git: failed to create gzip reader: %s", err)
			return
		}
		defer reader.Close()
	}

	svc := service.ServiceCommand{
		Dir:    repo.path,
		Stdout: w,
		Stdin:  reader,
	}

	if err := svc.UploadPack(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("failed to execute git-upload-pack: %s", err)
		return
	}
}
