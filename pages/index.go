// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package pages

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"

	gogit "github.com/go-git/go-git/v5"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (pages *Pages) index(w http.ResponseWriter, r *http.Request) {
	dirs, err := os.ReadDir(pages.core.ScanDir.Name())
	if err != nil {
		errors.WriteInternalServerError(pages.core, w, r)
		log.Printf("read scan path error: %s", err)
		return
	}

	summaries := make([]templates.RepositorySummary, 0, len(dirs))

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		path, err := pages.core.RepositoryPath(dir.Name())
		if err != nil {
			continue
		}

		name := filepath.Base(path)
		if slices.Contains(pages.core.Config.Repo.Unlisted, name) {
			continue
		}

		repo, err := gogit.PlainOpen(path)
		if err != nil {
			// This can be user having non-repository directory inside repo.scanPath.
			// In that case, logging an error would flood log output.
			// This makes debugging "repository not visible on top page" more difficult,
			// but is safer and steadier.
			continue
		}

		head, err := repo.Head()
		if err != nil {
			errors.WriteInternalServerError(pages.core, w, r)
			log.Printf("unable to get HEAD: %s", err)
			return
		}

		commit, err := repo.CommitObject(head.Hash())
		if err != nil {
			errors.WriteInternalServerError(pages.core, w, r)
			log.Printf("latest commit read error: %s", err)
			return
		}

		var category string
		if pages.core.Config.UI.Category.Grouping {
			category = git.GitwebCategory(repo)
		}

		summaries = append(summaries, templates.RepositorySummary{
			DisplayName: pages.core.RepositoryName(path),
			DirName:     name,
			Description: git.GitwebDescription(repo),
			Category:    category,
			LastCommit:  commit,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[j].LastCommit.Committer.When.Before(summaries[i].LastCommit.Committer.When)
	})

	data := templates.RepoListData{
		Config:       pages.core.Config,
		Repositories: summaries,
	}

	if err := pages.core.Template().ExecuteTemplate(w, "repo-list", data); err != nil {
		log.Println(err)
		return
	}
}
