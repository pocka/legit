// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package tests

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func SignatureAlice() *object.Signature {
	return &object.Signature{
		Name:  "Alice",
		Email: "alice@example.com",
		When:  time.Unix(0, 0),
	}
}

func CreateRepository(baseDir string, name string) (*git.Repository, *git.Worktree, error) {
	storage := filesystem.NewStorage(
		osfs.New(filepath.Join(baseDir, name), osfs.WithBoundOS()),
		cache.NewObjectLRUDefault(),
	)

	repo, err := git.InitWithOptions(storage, storage.Filesystem(), git.InitOptions{
		DefaultBranch: "refs/heads/trunk",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to create Git repository: %s", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to get Git worktree: %s", err)
	}

	return repo, wt, nil
}

func CreateBare(baseDir string, name string) error {
	storage := filesystem.NewStorage(
		osfs.New(filepath.Join(baseDir, fmt.Sprintf("%s.git", name)), osfs.WithBoundOS()),
		cache.NewObjectLRUDefault(),
	)

	target := url.URL{
		Scheme: "file",
		Path:   filepath.Join(baseDir, name),
	}

	_, err := git.Clone(storage, nil, &git.CloneOptions{
		URL:           target.String(),
		ReferenceName: "refs/heads/trunk",
	})

	return err
}
