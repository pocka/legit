// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) tree(pathRemainings string, w http.ResponseWriter, r *http.Request) {
	rev, ref, childPath, err := takeRef(repo.r, pathRemainings)
	if err != nil {
		errors.WriteNotFound(repo.core, w, r)
		log.Println(err)
		return
	}

	commit, err := repo.r.CommitObject(*rev)
	if err != nil {
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

	niceTree, err := git.NewNiceTree(tree, childPath)
	if err != nil {
		if err == object.ErrEntryNotFound {
			errors.WriteNotFound(repo.core, w, r)
			return
		}

		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	paths := []string{}
	if childPath != "" {
		paths = strings.Split(childPath, "/")
	}

	data := templates.RepoTreeRefData{
		Config: repo.core.Config,
		Meta: templates.RepositoryMeta{
			DisplayName: repo.core.RepositoryName(repo.dirname),
			DirName:     repo.dirname,
			Description: git.GitwebDescription(repo.r),
			Ref:         ref,
		},
		Path:  paths,
		Files: niceTree,
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-tree-ref", data); err != nil {
		log.Println(err)
		return
	}
}
