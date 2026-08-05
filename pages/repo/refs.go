// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"log"
	"net/http"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) refs(w http.ResponseWriter, r *http.Request) {
	tags, err := git.Tags(repo.r)
	if err != nil {
		// Non-fatal, we *should* have at least one branch to show.
		log.Println(err)
	}

	branches, err := git.Branches(repo.r)
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	mainBranch, err := git.FindMainBranch(repo.r, repo.core.Config.Repo.MainBranch)
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	data := templates.RepoRefsData{
		Config: repo.core.Config,
		Meta: templates.RepositoryMeta{
			DisplayName: repo.core.RepositoryName(repo.dirname),
			DirName:     repo.dirname,
			Description: git.GitwebDescription(repo.r),
			Ref:         mainBranch,
		},
		Tags:     tags,
		Branches: branches,
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-refs", data); err != nil {
		log.Println(err)
		return
	}
}
