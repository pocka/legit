// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package git

import (
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitsOptions struct {
	// Limit is maximum number of commits to retrieve.
	// 0 means no limit.
	Limit uint32

	Before plumbing.Hash
	After  plumbing.Hash
}

type CommitsPageMeta struct {
	HasNextPage bool
	HasPrevPage bool
}

func Commits(repo *git.Repository, rev plumbing.Hash, opts CommitsOptions) ([]*object.Commit, CommitsPageMeta, error) {
	if !opts.After.IsZero() {
		return commitsAfter(repo, rev, opts)
	} else if !opts.Before.IsZero() {
		return commitsBefore(repo, opts)
	} else {
		return commitsFromTip(repo, rev, opts)
	}
}

func commitsFromTip(repo *git.Repository, rev plumbing.Hash, opts CommitsOptions) ([]*object.Commit, CommitsPageMeta, error) {
	ci, err := repo.Log(&git.LogOptions{From: rev})
	if err != nil {
		return nil, CommitsPageMeta{}, fmt.Errorf("commits from ref: %w", err)
	}

	commits := make([]*object.Commit, 0, opts.Limit)
	meta := CommitsPageMeta{
		HasPrevPage: false,
		HasNextPage: true,
	}

	for range opts.Limit {
		if next, err := ci.Next(); err != nil {
			if err == io.EOF {
				meta.HasNextPage = false
				break
			} else {
				return nil, CommitsPageMeta{}, fmt.Errorf("log failure: %w", err)
			}
		} else {
			commits = append(commits, next)
		}
	}

	// For when the number of commits equals to "opts.Limit".
	if _, err := ci.Next(); err == io.EOF {
		meta.HasNextPage = false
	}

	return commits, meta, nil
}

func commitsBefore(repo *git.Repository, opts CommitsOptions) ([]*object.Commit, CommitsPageMeta, error) {
	ci, err := repo.Log(&git.LogOptions{
		From: opts.Before,
	})
	if err != nil {
		return nil, CommitsPageMeta{}, fmt.Errorf("commits from ref: %w", err)
	}

	commits := make([]*object.Commit, 0, opts.Limit)
	meta := CommitsPageMeta{
		HasPrevPage: true,
		HasNextPage: true,
	}

	iteration := opts.Limit + 1
	for range iteration {
		if next, err := ci.Next(); err != nil {
			if err == io.EOF {
				meta.HasNextPage = false
				break
			} else {
				return nil, CommitsPageMeta{}, fmt.Errorf("log failure: %w", err)
			}
		} else {
			if next.Hash != opts.Before {
				commits = append(commits, next)
			}
		}
	}

	return commits, meta, nil
}

func commitsAfter(repo *git.Repository, rev plumbing.Hash, opts CommitsOptions) ([]*object.Commit, CommitsPageMeta, error) {
	ci, err := repo.Log(&git.LogOptions{From: rev})
	if err != nil {
		return nil, CommitsPageMeta{}, fmt.Errorf("commits from ref: %w", err)
	}

	meta := CommitsPageMeta{
		HasPrevPage: false,
		HasNextPage: true,
	}

	queue := newRingBuffer[*object.Commit](opts.Limit)

	for {
		if next, err := ci.Next(); err != nil {
			if err == io.EOF {
				meta.HasNextPage = false
				break
			} else {
				return nil, CommitsPageMeta{}, fmt.Errorf("log failure: %w", err)
			}
		} else {
			if next.Hash == opts.After {
				break
			}

			queue.push(next)
		}
	}

	commits := queue.toSlice()

	if len(commits) > 0 {
		meta.HasPrevPage = rev != commits[0].Hash
	}

	return commits, meta, nil
}
