// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	htmlout "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"

	"github.com/pocka/legit/renderer/html"
)

type repoLinkTransformer struct {
	repo *Repo
	ref  string
}

func newRepoLinkTransformer(repo *Repo, ref string) *repoLinkTransformer {
	return &repoLinkTransformer{
		repo: repo,
		ref:  ref,
	}
}

func (t *repoLinkTransformer) RewriteInternalMediaSource(src string) string {
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
	basePath := fmt.Sprintf("/%s/blob/%s", t.repo.dirname, t.ref)
	path := filepath.Join(basePath, href)
	if strings.Index(path, basePath) != 0 {
		path = basePath + path
	}

	return path + query
}

func (t *repoLinkTransformer) RewriteInternalLink(link string) string {
	href := link

	if strings.IndexByte(href, '/') == 0 {
		// Repository index page is the "root" for repository's README file.
		href = "." + href
	}

	// Output path should not go beyond this.
	var basePath string
	if strings.LastIndexByte(href, '/') == len(href)-1 {
		basePath = fmt.Sprintf("/%s/tree/%s", t.repo.dirname, t.ref)
	} else {
		basePath = fmt.Sprintf("/%s/blob/%s", t.repo.dirname, t.ref)
	}
	path := filepath.Join(basePath, href)
	if strings.Index(path, basePath) != 0 {
		path = basePath + path
	}

	if t.repo.htmlRenderer(path) != nil && !strings.ContainsRune(path, '?') {
		path += "?preview=html"
	}

	return path
}

func (repo *Repo) htmlRenderer(filename string) html.Renderer {
	switch filepath.Ext(filename) {
	case ".md", ".mkd", ".markdown":
		return &repo.core.Markdown
	default:
		return nil
	}
}

func highlightCode(fileName string, code string) (template.HTML, error) {
	lexer := lexers.Get(fileName)
	if lexer == nil {
		return "", nil
	}

	formatter := htmlout.New(htmlout.WithClasses(true), htmlout.ClassPrefix("chroma-"))

	iter, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", fmt.Errorf("failed to tokenize code: %s", err)
	}

	var output bytes.Buffer
	err = formatter.Format(&output, styles.Fallback, iter)
	if err != nil {
		return "", fmt.Errorf("failed to highlight code: %s", err)
	}

	return template.HTML(output.String()), nil
}
