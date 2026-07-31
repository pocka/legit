// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pocka/legit/tests"
)

func TestGitwebDescriptionReadsDescriptionFile(t *testing.T) {
	root := t.TempDir()

	// non-bare
	{
		repo, worktree, err := tests.CreateRepository(root, "foo")
		if err != nil {
			t.Fatal(err)
		}

		_, err = worktree.Commit("init", &git.CommitOptions{
			AllowEmptyCommits: true,
			Author:            tests.SignatureAlice(),
		})
		if err != nil {
			t.Fatal(err)
		}

		dotgit := repo.Storer.(*filesystem.Storage)
		file, err := dotgit.Filesystem().Create("description")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := file.Write([]byte("Foo Bar")); err != nil {
			t.Fatal(err)
		}

		file.Close()

		r, err := Open(filepath.Join(root, "foo"), "trunk")
		if err != nil {
			t.Fatal(err)
		}

		description := r.GitwebDescription()
		if description != "Foo Bar" {
			t.Errorf("Unexpected description, got: %s", description)
		}
	}

	// bare
	{
		if err := tests.CreateBare(root, "foo"); err != nil {
			t.Fatal(err)
		}

		file, err := os.Create(filepath.Join(root, "foo.git", "description"))
		if err != nil {
			t.Fatal(err)
		}

		if _, err := file.WriteString("Foo Bare"); err != nil {
			t.Fatal(err)
		}

		file.Close()

		r, err := Open(filepath.Join(root, "foo.git"), "trunk")
		if err != nil {
			t.Fatal(err)
		}

		description := r.GitwebDescription()
		if description != "Foo Bare" {
			t.Errorf("Unexpected description, got: %s", description)
		}
	}
}

func TestGitwebDescriptionShouldNotAccessWorktree(t *testing.T) {
	root := t.TempDir()

	_, worktree, err := tests.CreateRepository(root, "foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = worktree.Commit("init", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author:            tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatal(err)
	}

	file, err := worktree.Filesystem.Create("description")
	if _, err := file.Write([]byte("Foo Bar")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	r, err := Open(filepath.Join(root, "foo"), "trunk")
	if err != nil {
		t.Fatal(err)
	}

	description := r.GitwebDescription()
	if description != "" {
		t.Errorf("Expected empty description, got: %s", description)
	}
}

func TestGitwebDescriptionReadsGitConfig(t *testing.T) {
	root := t.TempDir()

	repo, worktree, err := tests.CreateRepository(root, "foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = worktree.Commit("init", &git.CommitOptions{
		AllowEmptyCommits: true,
		Author:            tests.SignatureAlice(),
	})
	if err != nil {
		t.Fatal(err)
	}

	config, err := repo.Config()
	if err != nil {
		t.Fatal(err)
	}

	section := config.Raw.Section("gitweb")
	section.AddOption("description", "Foo Bar")
	if err := repo.SetConfig(config); err != nil {
		t.Fatal(err)
	}

	r, err := Open(filepath.Join(root, "foo"), "trunk")
	if err != nil {
		t.Fatal(err)
	}

	description := r.GitwebDescription()
	if description != "Foo Bar" {
		t.Errorf("Unexpected description, got: %s", description)
	}
}
