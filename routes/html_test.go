// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"testing"
)

func TestRewriteInternalMediaSource(t *testing.T) {
	r := NewRepoLinkTransformer("foo", "trunk")

	{
		input := "media/screenshot.png"
		expected := "/foo/blob/trunk/media/screenshot.png?raw"

		if out := r.RewriteInternalMediaSource(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "./media/screenshot.png"
		expected := "/foo/blob/trunk/media/screenshot.png?raw"

		if out := r.RewriteInternalMediaSource(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "/media/screenshot.png"
		expected := "/foo/blob/trunk/media/screenshot.png?raw"

		if out := r.RewriteInternalMediaSource(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "media/screenshot.png?width=500&height=300"
		expected := "/foo/blob/trunk/media/screenshot.png?width=500&height=300&raw"

		if out := r.RewriteInternalMediaSource(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}
}
