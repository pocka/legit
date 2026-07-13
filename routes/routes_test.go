// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/tests"
)

func TestGetSummaryOK(t *testing.T) {
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

	server := httptest.NewServer(Handlers(&c, embed.StaticDir(), embed.TemplatesDir()))
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
