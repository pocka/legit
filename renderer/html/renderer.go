// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package html

type Renderer interface {
	// Render processes "code" then returns HTML text.
	Render(code []byte, transformer Transformer) ([]byte, error)
}

type Transformer interface {
	// RewriteInternalMediaSource takes a site internal link to a media file
	// then returns a new link.
	RewriteInternalMediaSource(src string) string

	// RewriteInternalLink takes a site internal link then returns a new link.
	RewriteInternalLink(href string) string
}
