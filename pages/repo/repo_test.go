// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package repo

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

func FuzzRepoPages(f *testing.F) {
	repos := f.TempDir()

	secret, err := os.Create(filepath.Join(repos, "secret.txt"))
	if err != nil {
		f.Fatal(err)
	}
	defer secret.Close()

	if _, err := secret.WriteString("pepperoni"); err != nil {
		f.Fatal(err)
	}

	_, worktree, err := tests.CreateRepository(repos, "foo")
	if err != nil {
		f.Fatal(err)
	}

	readme, err := worktree.Filesystem.Create("README.md")
	if err != nil {
		f.Fatal(err)
	}

	if _, err := readme.Write([]byte("* iawsoiwjfngbhfg812uhjikwe6789asfd")); err != nil {
		f.Fatal(err)
	}

	_ = readme.Close()

	if _, err := worktree.Add("README.md"); err != nil {
		f.Fatal(err)
	}

	_, err = worktree.Commit("Add README", &git.CommitOptions{
		Author: tests.SignatureAlice(),
	})
	if err != nil {
		f.Fatalf("Unable to commit: %s", err)
	}

	var c config.Config
	c.Repo.ScanPath = repos
	c.Repo.Readme = []string{"README.md"}
	c.Repo.MainBranch = []string{"trunk"}

	core, err := core.New(&c)
	if err != nil {
		f.Fatal(err)
	}

	f.Add("/")
	f.Fuzz(func(t *testing.T, path string) {
		target := "/" + path
		url, err := url.Parse(target)
		if err != nil {
			// Skip non-GET-table URLs
			t.SkipNow()
		}

		// httptest.NewRequest panics even on slight erroneous path
		r, err := http.NewRequest(http.MethodGet, url.String(), nil)
		if err != nil {
			t.SkipNow()
		}

		handler, err := New(filepath.Join(repos, "foo"), r.URL.Path, core)
		if err != nil {
			t.Fatal(err)
		}

		r.URL.Path = "/foo" + r.URL.Path

		w := httptest.NewRecorder()

		// httptest.Server actually uses socket and if we run that, OS stops this
		// test due to socket usage limit.
		handler.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()

		// No InternalServerError
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusNotFound:
		default:
			t.Errorf(
				"Expected HTTP %d or %d, got %d",
				http.StatusOK, http.StatusNotFound,
				res.StatusCode,
			)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		body := string(b)

		if strings.Contains(body, "pepperoni") {
			t.Error("Returned body contains content of secret file under repos directory")
		}

		if !strings.Contains(body, "href=\"/static/noscript.css?r=") {
			t.Errorf("Returned HTML does not use an embedded template")
		}
	})
}
