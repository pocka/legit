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

	return fmt.Sprintf("/%s/blob/%s/%s?raw", t.repoName, t.ref, href)
}

func (t *RepoLinkTransformer) RewriteInternalLink(link string) string {
	href := link

	if strings.IndexByte(href, '/') == 0 {
		// Repository index page is the "root" for repository's README file.
		href = "." + href
	}

	if strings.LastIndexByte(href, '/') == len(href)-1 {
		return fmt.Sprintf("/%s/tree/%s/%s", t.repoName, t.ref, href[:len(href)-1])
	} else {
		return fmt.Sprintf("/%s/blob/%s/%s", t.repoName, t.ref, href)
	}
}

func (d *deps) htmlRenderer(file *object.File) html.Renderer {
	switch filepath.Ext(file.Name) {
	case ".md", ".mkd", ".markdown":
		return &d.markdown
	default:
		return nil
	}
}
