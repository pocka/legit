// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pocka/legit/renderer/html"
)

type RepoLinkTransformer struct {
	repoName string
	ref      string
}

func NewRepoLinkTransformer(repoName string, ref string) *RepoLinkTransformer {
	return &RepoLinkTransformer{
		repoName: repoName,
		ref:      ref,
	}
}

func (t *RepoLinkTransformer) RewriteInternalMediaSource(src string) string {
	href := src

	if strings.IndexByte(href, '/') == 0 {
		// Repository index page is the "root" for repository's README file.
		href = "." + href
	}

	query := "?raw"
	if strings.ContainsRune(href, '?') {
		query = "&raw"
	}

	// Output path should not go beyond this.
	basePath := fmt.Sprintf("/%s/blob/%s", t.repoName, t.ref)
	path := filepath.Join(basePath, href)
	if strings.Index(path, basePath) != 0 {
		path = basePath + path
	}

	return path + query
}

func (t *RepoLinkTransformer) RewriteInternalLink(link string) string {
	href := link

	if strings.IndexByte(href, '/') == 0 {
		// Repository index page is the "root" for repository's README file.
		href = "." + href
	}

	// Output path should not go beyond this.
	var basePath string
	if strings.LastIndexByte(href, '/') == len(href)-1 {
		basePath = fmt.Sprintf("/%s/tree/%s", t.repoName, t.ref)
	} else {
		basePath = fmt.Sprintf("/%s/blob/%s", t.repoName, t.ref)
	}
	path := filepath.Join(basePath, href)
	if strings.Index(path, basePath) != 0 {
		path = basePath + path
	}

	return path
}

func (d *deps) htmlRenderer(file *object.File) html.Renderer {
	switch filepath.Ext(file.Name) {
	case ".md", ".mkd", ".markdown":
		return &d.markdown
	default:
		return nil
	}
}
