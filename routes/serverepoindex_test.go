// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/tests"
)

func TestServeRepoIndexOK(t *testing.T) {
	repos := t.TempDir()

	_, worktree, err := tests.CreateRepository(repos, "foo")
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

	_, err = worktree.Commit("Add README", &git.CommitOptions{
		Author: tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatalf("Unable to commit: %s", err)
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	server := httptest.NewServer(Handler(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	target, err := url.JoinPath(server.URL, "/foo")
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get(target)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected HTTP 200, Got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(body), "<li>iawsoiwjfngbhfg812uhjikwe6789asfd</li>") {
		t.Error("Body not containing magic string")
	}
}

// TestServeRepoIndexPreventPathTraversal tests prevention of path traversal.
// https://github.com/icyphox/legit/issues/49
func TestServeRepoIndexPreventPathTraversal(t *testing.T) {
	root := t.TempDir()

	child := filepath.Join(root, "child")
	if err := os.Mkdir(child, 0o750); err != nil {
		t.Fatal(err)
	}

	{
		_, worktree, err := tests.CreateRepository(root, "private")
		if err != nil {
			t.Fatal(err)
		}

		readme, err := worktree.Filesystem.Create("README.md")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := readme.Write([]byte("MY SECRET!")); err != nil {
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
	}

	{
		_, worktree, err := tests.CreateRepository(child, "public")
		if err != nil {
			t.Fatal(err)
		}

		readme, err := worktree.Filesystem.Create("README.md")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := readme.Write([]byte("everybody can see this")); err != nil {
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
	}

	var c config.Config
	c.Repo.ScanPath = child
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	server := httptest.NewServer(Handler(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	target, err := url.JoinPath(server.URL, "..%2Fprivate")
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get(target)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected HTTP %d, Got %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestServeRepoIndexPreventIgnoredRepoReveal(t *testing.T) {
	repos := t.TempDir()

	{
		_, worktree, err := tests.CreateRepository(repos, "private")
		if err != nil {
			t.Fatal(err)
		}

		readme, err := worktree.Filesystem.Create("README.md")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := readme.Write([]byte("MY SECRET!")); err != nil {
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
	}

	{
		_, worktree, err := tests.CreateRepository(repos, "public")
		if err != nil {
			t.Fatal(err)
		}

		readme, err := worktree.Filesystem.Create("README.md")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := readme.Write([]byte("everybody can see this")); err != nil {
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
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}
	c.Repo.Ignore = []string{"private"}

	server := httptest.NewServer(Handler(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	target, err := url.JoinPath(server.URL, ".%2Fprivate")
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get(target)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected HTTP %d, Got %d", http.StatusNotFound, res.StatusCode)
	}
}
