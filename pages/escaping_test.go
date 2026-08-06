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
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/core"
	"github.com/pocka/legit/tests"
)

func TestRepoTreeHandleSlashes(t *testing.T) {
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

	core, err := core.New(&c)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(New(core))
	defer server.Close()

	for _, p := range []string{"/", "/./", "//", "/////////////////"} {
		target, err := url.JoinPath(server.URL, "/foo/tree/trunk")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target + p)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("(%s) Expected HTTP %d, Got %d", p, http.StatusOK, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "href=\"/foo/blob/trunk/README.md\"") {
			t.Errorf("(%s) No link for README.md", p)
			t.Log(string(body))
		}
	}
}

func TestRepoTreePreventsPathTraversal(t *testing.T) {
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

	core, err := core.New(&c)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(New(core))
	defer server.Close()

	for _, p := range []string{"/../", "/.././.."} {
		target, err := url.JoinPath(server.URL, "/foo/tree/trunk")
		if err != nil {
			t.Fatal(err)
		}

		res, err := http.Get(target + p)
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("(%s) Expected HTTP %d, Got %d", p, http.StatusNotFound, res.StatusCode)
		}
	}
}
