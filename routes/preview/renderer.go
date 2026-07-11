// This file contains preview renderers.
//
// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package preview

import (
	"bytes"
	"net/url"
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
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))

	tree := parser.Parse(code)
	tree.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		if node.Type == blackfriday.Image {
			src := string(node.LinkData.Destination)
			parsedURL, err := url.Parse(src)
			if err == nil && parsedURL.Host != "" {
				// External URL, skip.
				return blackfriday.GoToNext
			}

			queries := parsedURL.Query()
			queries.Add("raw", "1")
			parsedURL.RawQuery = queries.Encode()

			node.LinkData.Destination = []byte(parsedURL.String())
		}

		return blackfriday.GoToNext
	})

	var writer bytes.Buffer
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	renderer.RenderHeader(&writer, tree)
	tree.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(&writer, node, entering)
	})
	renderer.RenderFooter(&writer, tree)

	unsafe := blackfriday.Run(writer.Bytes(), blackfriday.WithExtensions(blackfriday.CommonExtensions))
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
