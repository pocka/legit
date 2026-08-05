// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"log"
	"net/http"

	"github.com/pocka/legit/git/service"
)

func (repo *Repo) inforefs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")
	w.WriteHeader(http.StatusOK)

	svc := service.ServiceCommand{
		Dir:    repo.path,
		Stdout: w,
	}

	if err := svc.InfoRefs(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("failed to execute git-upload-pack (info/refs): %s", err)
		return
	}
}
