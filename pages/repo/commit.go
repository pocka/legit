// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"encoding/hex"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) commit(pathRemainings string, w http.ResponseWriter, r *http.Request) {
	h, err := hex.DecodeString(pathRemainings)
	if err != nil {
		errors.WriteNotFound(repo.core, w, r)
		return
	}

	var hash plumbing.Hash
	copy(hash[:], h)

	commit, err := repo.r.CommitObject(hash)
	if err != nil {
		if err == plumbing.ErrObjectNotFound {
			errors.WriteNotFound(repo.core, w, r)
			return
		}

		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	var diff *git.NiceDiff
	if repo.core.Config.Repo.Diff == "system" {
		diff, err = git.SystemDiff(repo.path, commit)
	} else {
		diff, err = git.Diff(commit)
	}
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	data := templates.RepoCommitData{
		Config: repo.core.Config,
		Meta: templates.RepositoryMeta{
			DisplayName: repo.core.RepositoryName(repo.dirname),
			DirName:     repo.dirname,
			Description: git.GitwebDescription(repo.r),
			Ref:         hash.String(),
		},
		Commit: commit,
		Parent: diff.Parent,
		Diff:   diff,
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-commit", data); err != nil {
		log.Println(err)
		return
	}
}
