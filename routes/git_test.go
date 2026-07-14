// This file tests git operations.
//
// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/tests"
)

func TestGitCloneOK(t *testing.T) {
	repos := t.TempDir()

	_, worktree, err := tests.CreateRepository(repos, "foo")
	if err != nil {
		t.Fatal(err)
	}

	readme, err := worktree.Filesystem.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := readme.Write([]byte("# Foo")); err != nil {
		t.Fatal(err)
	}

	_ = readme.Close()

	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatal(err)
	}

	hash, err := worktree.Commit("Add README", &git.CommitOptions{
		Author: tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatalf("Unable to commit: %s", err)
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	server := httptest.NewServer(Handlers(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	cloneURL, err := url.JoinPath(server.URL, "/foo")
	if err != nil {
		t.Fatal(err)
	}

	cloned, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           cloneURL,
		ReferenceName: "refs/heads/trunk",
	})
	if err != nil {
		t.Fatal(err)
	}

	commit, err := cloned.CommitObject(hash)
	if err != nil {
		t.Fatalf("Unable to get commit from the cloned repository: %s", err)
	}

	if commit.Message != "Add README" {
		t.Fatalf("Unexpected commit message (%s)", commit.Message)
	}
}

func TestGitBareCloneOK(t *testing.T) {
	repos := t.TempDir()

	_, worktree, err := tests.CreateRepository(repos, "foo")
	if err != nil {
		t.Fatal(err)
	}

	readme, err := worktree.Filesystem.Create("README.md")
	if err != nil {
		t.Fatal(err)
	}

	readmeContent := `# Foo

This is foo.
I'm [Markdown](https://commonmark.org/) file for *load* **testing**.
`

	if _, err := readme.Write([]byte(readmeContent)); err != nil {
		t.Fatal(err)
	}

	_ = readme.Close()

	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatal(err)
	}

	hash, err := worktree.Commit("Add README", &git.CommitOptions{
		Author: tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatalf("Unable to commit: %s", err)
	}

	if err := tests.CreateBare(repos, "foo"); err != nil {
		t.Fatalf("Unable to create bare repository: %s", err)
	}

	if err := os.RemoveAll(filepath.Join(repos, "foo")); err != nil {
		t.Fatal(err)
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	server := httptest.NewServer(Handlers(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	cloneURL, err := url.JoinPath(server.URL, "/foo.git")
	if err != nil {
		t.Fatal(err)
	}

	cloned, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           cloneURL,
		ReferenceName: "refs/heads/trunk",
	})
	if err != nil {
		t.Fatal(err)
	}

	commit, err := cloned.CommitObject(hash)
	if err != nil {
		t.Fatalf("Unable to get commit from the cloned repository: %s", err)
	}

	if commit.Message != "Add README" {
		t.Fatalf("Unexpected commit message (%s)", commit.Message)
	}
}

// https://github.com/icyphox/legit/issues/56
// https://tangled.org/pocka.jp/legit/issues/2
func TestCloneRespectVisibilityConfig(t *testing.T) {
	repos := t.TempDir()

	for _, name := range []string{"public", "ignored", "unlisted"} {
		_, worktree, err := tests.CreateRepository(repos, name)
		if err != nil {
			t.Fatal(err)
		}

		readme, err := worktree.Filesystem.Create("README.md")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := readme.Write([]byte(name)); err != nil {
			t.Fatal(err)
		}

		_ = readme.Close()

		if _, err := worktree.Add("README.md"); err != nil {
			t.Fatal(err)
		}

		_, err = worktree.Commit("Add README", &git.CommitOptions{
			Author: tests.SignatureAlice(),
		})
		if err != nil {
			t.Fatalf("Unable to commit: %s", err)
		}

		if err := tests.CreateBare(repos, name); err != nil {
			t.Fatalf("Unable to create bare repository: %s", err)
		}

		if err := os.RemoveAll(filepath.Join(repos, name)); err != nil {
			t.Fatal(err)
		}
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}
	c.Repo.Ignore = []string{"ignored.git"}
	c.Repo.Unlisted = []string{"unlisted.git"}

	server := httptest.NewServer(Handlers(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	publicUrl, err := url.JoinPath(server.URL, "/public.git")
	if err != nil {
		t.Fatal(err)
	}

	_, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           publicUrl,
		ReferenceName: "refs/heads/trunk",
	})
	if err != nil {
		t.Error(err)
	}

	ignoredUrl, err := url.JoinPath(server.URL, "/ignored.git")
	if err != nil {
		t.Fatal(err)
	}

	_, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           ignoredUrl,
		ReferenceName: "refs/heads/trunk",
	})
	if err == nil {
		t.Error("Cloning ignored repository succeeded (expected an error.)")
	}

	unlistedUrl, err := url.JoinPath(server.URL, "/unlisted.git")
	if err != nil {
		t.Fatal(err)
	}

	_, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           unlistedUrl,
		ReferenceName: "refs/heads/trunk",
	})
	if err != nil {
		t.Error(err)
	}
}
