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
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/tests"
)

func TestIgnoreSymlinkByDefault(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	_, worktree, err := tests.CreateRepository(src, "foo")
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

	os.Symlink(filepath.Join(src, "foo"), filepath.Join(dest, "foo"))

	var c config.Config
	c.Repo.ScanPath = dest
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	server := httptest.NewServer(Handlers(&c, embed.StaticDir(), embed.TemplatesDir()))
	defer server.Close()

	topPage, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	topPageBody, err := io.ReadAll(topPage.Body)
	topPage.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(topPageBody), ">foo</") {
		t.Error("Top page lists symlink-ed repository.")
	}

	target, err := url.JoinPath(server.URL, "/foo")
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get(target)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected HTTP %d, Got %d", http.StatusNotFound, res.StatusCode)
	}

	_, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           target,
		ReferenceName: "refs/heads/trunk",
	})
	if err == nil {
		t.Error("Expected clone error, got successful clone.")
	}
}
