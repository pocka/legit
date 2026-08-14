// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) index(w http.ResponseWriter, r *http.Request) {
	mainBranch, err := git.FindMainBranch(repo.r, repo.core.Config.Repo.MainBranch)
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	rev, err := repo.r.ResolveRevision(plumbing.Revision(mainBranch))
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Printf("unable to fetch latest commit hash: %s", err)
		return
	}

	latestCommit, err := repo.r.CommitObject(*rev)
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Printf("unable to load latest commit: %s", err)
		return
	}

	commits, _, err := git.Commits(repo.r, *rev, git.CommitsOptions{Limit: 3})
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	tree, err := latestCommit.Tree()
	if err != nil {
		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	transformer := newRepoLinkTransformer(repo, mainBranch)

	var readmeContent template.HTML
	for _, readme := range repo.core.Config.Repo.Readme {
		file, _ := tree.File(readme)
		if file == nil {
			continue
		}

		content, _ := file.Contents()

		// Skip empty files.
		if len(content) > 0 {
			renderer := repo.htmlRenderer(file.Name)
			if renderer == nil {
				renderer = &repo.core.Plaintext
			}

			result, err := renderer.Render([]byte(content), transformer)
			if err != nil {
				log.Printf("Unable to render %s/%s, skipping.", repo.dirname, readme)
				continue
			}

			readmeContent = template.HTML(result)
			break
		}
	}

	isGoModule := false
	if _, err := tree.File("go.mod"); err == nil {
		isGoModule = true
	}

	data := templates.RepoTopData{
		Config: repo.core.Config,
		Meta: templates.RepositoryMeta{
			DisplayName: repo.core.RepositoryName(repo.dirname),
			DirName:     repo.dirname,
			Description: git.GitwebDescription(repo.r),
			Ref:         mainBranch,
		},
		Readme:        readmeContent,
		DefaultBranch: mainBranch,
		RecentCommits: commits,
		IsGoModule:    isGoModule,
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-top", data); err != nil {
		log.Println(err)
		return
	}
}
