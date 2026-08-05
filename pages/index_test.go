// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package pages

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
	"github.com/pocka/legit/core"
	"github.com/pocka/legit/tests"
)

func TestServeIndexIgnoreNonGitRepos(t *testing.T) {
	repos := t.TempDir()

	// Normal repository
	{
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
	}

	// Non git directory
	{
		if err := os.Mkdir(filepath.Join(repos, "bar"), 0o777); err != nil {
			t.Fatal(err)
		}
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

	target, err := url.JoinPath(server.URL, "/")
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

	if !strings.Contains(string(body), "foo") {
		t.Error("Body not containing magic string")
	}
}
