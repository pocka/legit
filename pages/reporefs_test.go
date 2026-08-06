// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package pages

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/core"
	"github.com/pocka/legit/tests"
)

// TestRepoHandleSlashRefsOK tests repository router handles slashed refs
// such as "feature/new-banana" correctly.
func TestRepoHandleSlashRefsOK(t *testing.T) {
	repos := t.TempDir()

	repo, worktree, err := tests.CreateRepository(repos, "foo")
	if err != nil {
		t.Fatal(err)
	}

	readme, err := worktree.Filesystem.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := readme.Write([]byte("* iawsoiwjfngbhfg812uhjikwe6789asfd")); err != nil {
		t.Fatal(err)
	}

	_ = readme.Close()

	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatal(err)
	}

	commit, err := worktree.Commit("Add README", &git.CommitOptions{
		Author: tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatalf("Unable to commit: %s", err)
	}

	if err := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("feature/john/new-renderer"),
		Create: true,
	}); err != nil {
		t.Fatalf("Unable to create a branch: %s", err)
	}

	if _, err := repo.CreateTag("foo/bar/baz", commit, nil); err != nil {
		t.Fatalf("Unable to create a tag: %s", err)
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	core, err := core.New(&c)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(New(core))
	defer server.Close()

	// List refs
	{
		target, err := url.JoinPath(server.URL, "/foo/refs")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected HTTP %d, Got %d", http.StatusOK, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "href=\"/foo/log/feature/john/new-renderer\"") {
			t.Error("No link for feature/john/new-renderer branch")
		}

		if !strings.Contains(string(body), "href=\"/foo/log/trunk\"") {
			t.Error("No link for trunk branch")
		}

		if !strings.Contains(string(body), "href=\"/foo/log/foo/bar/baz\"") {
			t.Error("No link for foo/bar/baz tag")
		}
	}

	// Browse files of "foo/bar/baz" tag
	{
		target, err := url.JoinPath(server.URL, "/foo/tree/foo/bar/baz")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected HTTP %d, Got %d", http.StatusOK, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "href=\"/foo/blob/foo/bar/baz/README.md\"") {
			t.Error("No link for foo/bar/baz/README.md")
		}
	}

	// View "README.md" at "foo/bar/baz" tag
	{
		target, err := url.JoinPath(server.URL, "/foo/blob/foo/bar/baz/README.md")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected HTTP %d, Got %d", http.StatusOK, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "* iawsoiwjfngbhfg812uhjikwe6789asfd") {
			t.Error("README.md does not containg magic string")
		}
	}

	// Preview "README.md" in HTML at "foo/bar/baz" tag
	{
		target, err := url.JoinPath(server.URL, "/foo/blob/foo/bar/baz/README.md")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target + "?preview=html")
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected HTTP %d, Got %d", http.StatusOK, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "<li>iawsoiwjfngbhfg812uhjikwe6789asfd</li>") {
			t.Error("HTML preview of README.md does not containg magic string")
		}
	}
}
