// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) log(pathRemainings string, w http.ResponseWriter, r *http.Request) {
	rev, ref, childPath, err := takeRef(repo.r, pathRemainings)
	if err != nil {
		errors.WriteNotFound(repo.core, w, r)
		log.Println(err)
		return
	}

	if childPath != "" {
		errors.WriteNotFound(repo.core, w, r)
		return
	}

	limit := repo.core.Config.UI.CommitsPageSize
	opts := git.CommitsOptions{Limit: limit}
	query := r.URL.Query()

	if after := query.Get("before"); plumbing.IsHash(after) {
		opts.Before = plumbing.NewHash(after)
	}
	if before := query.Get("after"); plumbing.IsHash(before) {
		opts.After = plumbing.NewHash(before)
	}

	commits, page, err := git.Commits(repo.r, *rev, opts)
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	prevPageHref := ""
	nextPageHref := ""

	if len(commits) > 0 {
		if page.HasPrevPage {
			prevPageHref = fmt.Sprintf("/%s/log/%s?after=%s", repo.dirname, ref, commits[0].Hash.String())
		}

		if page.HasNextPage {
			nextPageHref = fmt.Sprintf("/%s/log/%s?before=%s", repo.dirname, ref, commits[len(commits)-1].Hash.String())
		}
	}

	data := templates.RepoLogRefData{
		Config: repo.core.Config,
		Meta: templates.RepositoryMeta{
			DisplayName: repo.core.RepositoryName(repo.dirname),
			DirName:     repo.dirname,
			Description: git.GitwebDescription(repo.r),
			Ref:         ref,
		},
		Commits:      commits,
		PrevPageHref: prevPageHref,
		NextPageHref: nextPageHref,
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-log-ref", data); err != nil {
		log.Println(err)
		return
	}
}
