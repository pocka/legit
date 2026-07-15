// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package html

import (
	"bytes"
	"net/url"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// MarkdownRenderer is a renderer that takes Markdown code and emits HTML.
// Use "NewMarkdownRenderer" to instantiate.
type MarkdownRenderer struct {
	// policy is a sanitization policy.
	policy *bluemonday.Policy
}

func NewMarkdownRenderer(policy *bluemonday.Policy) MarkdownRenderer {
	p := policy
	if p == nil {
		p = bluemonday.UGCPolicy()
	}

	return MarkdownRenderer{
		policy: p,
	}
}

func (r *MarkdownRenderer) Render(code []byte, transformer Transformer) ([]byte, error) {
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	tree := parser.Parse(code)

	if transformer != nil {
		tree.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if !entering {
				return blackfriday.GoToNext
			}

			if node.Type == blackfriday.Link || node.Type == blackfriday.Image {
				href := string(node.LinkData.Destination)
				parsedURL, err := url.Parse(href)
				if err == nil && parsedURL.Host != "" {
					// Full URL, skip.
					return blackfriday.GoToNext
				}

				if node.Type == blackfriday.Image {
					node.LinkData.Destination = []byte(transformer.RewriteInternalMediaSource(href))
				} else {
					node.LinkData.Destination = []byte(transformer.RewriteInternalLink(href))
				}

			}

			return blackfriday.GoToNext
		})
	}

	var writer bytes.Buffer
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})
	renderer.RenderHeader(&writer, tree)
	tree.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(&writer, node, entering)
	})
	renderer.RenderFooter(&writer, tree)

	unsafe := blackfriday.Run(
		writer.Bytes(),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)
	return r.policy.SanitizeBytes(unsafe), nil
}
