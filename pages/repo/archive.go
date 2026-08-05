// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
)

const (
	tgzSuffix = ".tar.gz"
)

func (repo *Repo) archive(pathRemainings string, w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(pathRemainings, tgzSuffix) {
		errors.WriteNotFound(repo.core, w, r)
		return
	}

	rev, ref, childPath, err := takeRef(repo.r, strings.TrimSuffix(pathRemainings, tgzSuffix))
	if err != nil {
		errors.WriteNotFound(repo.core, w, r)
		log.Println(err)
		return
	}

	// e.g., "foo/bar/baz.tar.gz" where "foo/bar" is the right ref.
	if childPath != "" {
		errors.WriteNotFound(repo.core, w, r)
		return
	}

	commit, err := repo.r.CommitObject(*rev)
	if err != nil {
		if err == plumbing.ErrObjectNotFound {
			errors.WriteNotFound(repo.core, w, r)
			return
		}

		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	prefix := fmt.Sprintf("%s-%s", repo.core.RepositoryName(repo.dirname), ref)
	prefix = strings.ReplaceAll(prefix, "/", "-")

	header := w.Header()
	header.Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s%s\"", prefix, tgzSuffix))
	header.Add("Content-Type", "application/gzip")

	// We can't return error pages from this on.

	gz := gzip.NewWriter(w)
	defer gz.Close()

	if err := git.WriteTar(gz, prefix, tree); err != nil {
		log.Println(err)
		return
	}

	if err := gz.Flush(); err != nil {
		log.Println(err)
		return
	}
}
