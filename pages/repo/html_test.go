// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

import (
	"testing"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/core"
)

func TestRewriteInternalMediaSource(t *testing.T) {
	var c config.Config
	c.Repo.ScanPath = t.TempDir()
	core, err := core.New(&c)
	if err != nil {
		t.Fatal(err)
	}

	repo := Repo{
		dirname: "foo",
		core:    core,
	}

	r := newRepoLinkTransformer(&repo, "trunk", "bar.md")

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

	{
		input := "../../../../favicon.ico"
		expected := "/foo/blob/trunk/favicon.ico?raw"

		if out := r.RewriteInternalMediaSource(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}
}

func TestRewriteInternalLink(t *testing.T) {
	var c config.Config
	c.Repo.ScanPath = t.TempDir()
	core, err := core.New(&c)
	if err != nil {
		t.Fatal(err)
	}

	repo := Repo{
		dirname: "foo",
		core:    core,
	}

	r := newRepoLinkTransformer(&repo, "trunk", "bar.md")

	{
		input := "CHANGELOG.md"
		expected := "/foo/blob/trunk/CHANGELOG.md?preview=html"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "docs/DEVELOPMENT.adoc"
		expected := "/foo/blob/trunk/docs/DEVELOPMENT.adoc"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "vendor/libfoo/"
		expected := "/foo/tree/trunk/vendor/libfoo"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "./vendor/libfoo/"
		expected := "/foo/tree/trunk/vendor/libfoo"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "/vendor/"
		expected := "/foo/tree/trunk/vendor"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := ".././../././.././../foo/../.././"
		expected := "/foo/tree/trunk/"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}

	{
		input := "../../../../legit/config.yaml"
		expected := "/foo/blob/trunk/legit/config.yaml"

		if out := r.RewriteInternalLink(input); out != expected {
			t.Errorf("Expected %s, got %s", expected, out)
		}
	}
}
