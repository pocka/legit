// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package git

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pocka/legit/tests"
)

func createCommitsTestRepo(t *testing.T, numberOfCommits int) (*GitRepo, []plumbing.Hash) {
	repos := t.TempDir()
	_, worktree, err := tests.CreateRepository(repos, "foo")
	if err != nil {
		t.Fatal(err)
	}

	commits := make([]plumbing.Hash, numberOfCommits)

	for i := range numberOfCommits {
		file, err := worktree.Filesystem.Create(fmt.Sprintf("c%d.txt", i))
		if err != nil {
			t.Fatal(err)
		}

		if _, err := file.Write([]byte(".")); err != nil {
			t.Fatal(err)
		}

		_ = file.Close()

		if _, err := worktree.Add(fmt.Sprintf("c%d.txt", i)); err != nil {
			t.Fatal(err)
		}

		hash, err := worktree.Commit(fmt.Sprintf("c%d", i), &git.CommitOptions{
			Author: tests.SignatureAlice(),
		})
		if err != nil {
			t.Fatal(err)
		}

		commits[numberOfCommits-1-i] = hash
	}

	repo, err := Open(filepath.Join(repos, "foo"), "trunk")
	if err != nil {
		t.Fatal(err)
	}

	return repo, commits
}

func TestCommitsOK(t *testing.T) {
	// Reusing repo & commits as creating commits is not cheap.
	repo, commits := createCommitsTestRepo(t, 50)

	// len(commits) == opts.Limit
	{
		log, meta, err := repo.Commits(CommitsOptions{
			Limit: 50,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(log) != len(commits) {
			t.Fatalf("Expected %d commits, got %d", len(commits), len(log))
		}

		if meta.HasNextPage {
			t.Error("HasNextPage is true (should be false)")
		}

		if meta.HasPrevPage {
			t.Error("HasPrevPage is true (should be false)")
		}

		for i, commit := range log {
			if commit.Hash != commits[i] {
				t.Errorf(
					"%dth commit has unexpected hash: expected %s, got %s",
					i, commits[i].String(), commit.Hash.String(),
				)
			}
		}
	}

	// len(commits) < opts.Limit
	{
		log, meta, err := repo.Commits(CommitsOptions{
			Limit: 100,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(log) != len(commits) {
			t.Fatalf("Expected %d commits, got %d", len(commits), len(log))
		}

		if meta.HasNextPage {
			t.Error("HasNextPage is true (should be false)")
		}

		if meta.HasPrevPage {
			t.Error("HasPrevPage is true (should be false)")
		}

		for i, commit := range log {
			if commit.Hash != commits[i] {
				t.Errorf(
					"%dth commit has unexpected hash: expected %s, got %s",
					i, commits[i].String(), commit.Hash.String(),
				)
			}
		}
	}

	// len(commits) > opts.Limit
	{
		log, meta, err := repo.Commits(CommitsOptions{
			Limit: 10,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(log) != 10 {
			t.Fatalf("Expected 10 commits, got %d", len(log))
		}

		if !meta.HasNextPage {
			t.Error("HasNextPage is false (should be true)")
		}

		if meta.HasPrevPage {
			t.Error("HasPrevPage is true (should be false)")
		}

		for i, commit := range log {
			if commit.Hash != commits[i] {
				t.Errorf(
					"%dth commit has unexpected hash: expected %s, got %s",
					i, commits[i].String(), commit.Hash.String(),
				)
			}
		}
	}

	// len(commits) > opts.Limit && Before
	{
		log, meta, err := repo.Commits(CommitsOptions{
			Limit:  5,
			Before: commits[15],
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(log) != 5 {
			t.Fatalf("Expected 5 commits, got %d", len(log))
		}

		if !meta.HasNextPage {
			t.Error("HasNextPage is false (should be true)")
		}

		if !meta.HasPrevPage {
			t.Error("HasPrevPage is false (should be true)")
		}

		for i, commit := range log {
			if commit.Hash != commits[i+16] {
				t.Errorf(
					"%dth commit has unexpected hash: expected %s, got %s",
					i, commits[i].String(), commit.Hash.String(),
				)
			}
		}
	}

	// len(commits) > opts.Limit && After
	{
		log, meta, err := repo.Commits(CommitsOptions{
			Limit: 5,
			After: commits[40],
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(log) != 5 {
			t.Fatalf("Expected 5 commits, got %d", len(log))
		}

		if !meta.HasNextPage {
			t.Error("HasNextPage is false (should be true)")
		}

		if !meta.HasPrevPage {
			t.Error("HasPrevPage is false (should be true)")
		}

		for i, commit := range log {
			if commit.Hash != commits[35+i] {
				t.Errorf(
					"%dth commit has unexpected hash: expected %s, got %s",
					i, commits[i].String(), commit.Hash.String(),
				)
			}
		}
	}
}
