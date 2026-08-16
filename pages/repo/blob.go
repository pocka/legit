// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
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

		data, err := io.ReadAll(reader)
		if err != nil {
			errors.WriteInternalServerError(repo.core, w, r)
			log.Println(err)
			return
		}

		header := w.Header()

		switch filepath.Ext(childPath) {
		case ".svg", ".SVG":
			// net/http.DetectContentType inappropriately uses MIME sniffing algorithm
			// for *user agents*. Because the spec intentionally disallows sniffing to
			// SVG, we have to manually "detect" SVG files and add appropriate content
			// type.
			header.Set("Content-Type", "image/svg+xml")
			header.Set("X-Content-Type-Options", "nosniff")

			// Raw contents are completely user (committer) controlled, so use of inline
			// style is perfectly fine. SVG graphics often use <style> element for
			// essential styles. Remote resouces are still prohibited.
			header.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")

			w.WriteHeader(http.StatusOK)

			if _, err := w.Write(data); err != nil {
				log.Printf("SVG write failed: %s", err)
			}
			return
		}

		sniffed := http.DetectContentType(data)

		if strings.HasPrefix(sniffed, "image/") ||
			strings.HasPrefix(sniffed, "video/") ||
			strings.HasPrefix(sniffed, "audio/") ||
			sniffed == "application/octet-stream" {
			// Serve media file and binary file as-is, and let browsers sniff if it wants to.
			// The reason this branch allows sniffing is the "DetectContentType" hardcodes
			// client MIME sniffing alghorithm and that can be outdated due to Go standard
			// library lagged behind browsers or use of old Go version for compiling legit.
			// As both server (this code) and client (browser) use the same spec, the risk of
			// deviation should be low enough.
			header.Set("Content-Type", sniffed)

			// In case DetectContentType detects DOM-enabled content (image/svg+xml).
			// As of Go v1.25, that function never returns "image/svg+xml" though.
			// The only effect this directive does to non-HTML/SVG files is it prevents
			// page (browser tab) from loading a favicon, which is fine to this endpoint.
			header.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")

			w.WriteHeader(http.StatusOK)

			if _, err := w.Write(data); err != nil {
				log.Printf("media or binary write failed: %s", err)
			}
			return
		}

		// Everything else should be served as plain text, to prevent browsers from
		// automatically doing their magics. The primary purpose of this is to prevent
		// attacker from distributing malicious HTML/JS/CSS files via this endpoint.
		header.Set("Content-Type", "text/plain")
		header.Set("X-Content-Type-Options", "nosniff")

		// In case browser ignored nosniff directive and opens the content as HTML/SVG.
		header.Set("Content-Security-Policy", "default-src 'none'")

		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(data); err != nil {
			log.Printf("write failed: %s", err)
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
