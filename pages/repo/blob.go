// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pocka/legit/git"
	"github.com/pocka/legit/pages/errors"
	"github.com/pocka/legit/templates"
)

func (repo *Repo) blob(pathRemainings string, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

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

	file, err := tree.File(childPath)
	if err != nil {
		if err == object.ErrFileNotFound {
			errors.WriteNotFound(repo.core, w, r)
			return
		}

		errors.WriteInternalServerError(repo.core, w, r)
		log.Println(err)
		return
	}

	if query.Has("raw") {
		reader, err := file.Reader()
		if err != nil {
			errors.WriteInternalServerError(repo.core, w, r)
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")

		if _, err := io.Copy(w, reader); err != nil {
			log.Printf("io copy failed: %s", err)
		}

		return
	}

	var contents string
	if isBin, _ := file.IsBinary(); isBin {
		contents = "Not displaying binary file"
	} else {
		contents, err = file.Contents()
		if err != nil {
			errors.WriteInternalServerError(repo.core, w, r)
			fmt.Printf("reading file contents failed: %s", err)
			return
		}
	}

	paths := []string{}
	if childPath != "" {
		paths = strings.Split(childPath, "/")
	}

	meta := templates.RepositoryMeta{
		DisplayName: repo.core.RepositoryName(repo.dirname),
		DirName:     repo.dirname,
		Description: git.GitwebDescription(repo.r),
		Ref:         ref,
	}

	if query.Has("preview") {
		previewType := query.Get("preview")
		switch previewType {
		case "html":
			renderer := repo.htmlRenderer(file.Name)
			if renderer == nil {
				errors.WriteNotFound(repo.core, w, r)
				log.Printf("Requested HTML preview for %s/%s, but the filetype has no HTML renderer", repo.dirname, childPath)
				return
			}

			html, err := renderer.Render([]byte(contents), newRepoLinkTransformer(repo, ref, file.Name))
			if err != nil {
				errors.WriteInternalServerError(repo.core, w, r)
				log.Printf("Failed to render HTML preview: %s", err)
				return
			}

			data := templates.RepoBlobRefHTMLPreviewData{
				Config:  repo.core.Config,
				Meta:    meta,
				Path:    paths,
				Content: template.HTML(html),
			}

			if err := repo.core.Template().ExecuteTemplate(w, "repo-blob-ref-html-preview", data); err != nil {
				log.Println(err)
				return
			}

			return
		default:
			errors.WriteNotFound(repo.core, w, r)
			return
		}
	}

	lines := make([]uint, strings.Count(contents, "\n"))
	for i := range lines {
		if i < 0 {
			continue
		}

		lines[i] = uint(i + 1)
	}

	previewTypes := make([]string, 0, 1)
	if repo.htmlRenderer(file.Name) != nil {
		previewTypes = append(previewTypes, "html")
	}

	data := templates.RepoBlobRefData{
		Config:       repo.core.Config,
		Meta:         meta,
		Path:         paths,
		Content:      contents,
		LineNumbers:  lines,
		PreviewTypes: previewTypes,
	}

	if repo.core.Config.Meta.SyntaxHighlight {
		if highlighted, err := highlightCode(childPath, contents); err == nil {
			data.SyntaxHighlightedContent = highlighted
		} else {
			log.Println(err)
		}
	}

	if err := repo.core.Template().ExecuteTemplate(w, "repo-blob-ref", data); err != nil {
		log.Println(err)
		return
	}
}
