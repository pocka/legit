// This file contains preview renderers.
//
// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package preview

import (
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type Renderer interface {
	GetPreviewType() string

	Render(code []byte) ([]byte, error)
}

type MarkdownToHtmlRenderer struct{}

func (r MarkdownToHtmlRenderer) GetPreviewType() string {
	return "html"
}

func (r MarkdownToHtmlRenderer) Render(code []byte) ([]byte, error) {
	sanitizer := bluemonday.UGCPolicy()
	unsafe := blackfriday.Run(code, blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return sanitizer.SanitizeBytes(unsafe), nil
}

func GetPreviewRenderers(fileName string) []Renderer {
	ext := filepath.Ext(fileName)

	switch ext {
	case ".md", ".mkd", ".markdown":
		return []Renderer{MarkdownToHtmlRenderer{}}
	default:
		return []Renderer{}
	}
}
