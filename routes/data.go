// This file defines the data passed to templates.
//
// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"html/template"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pocka/legit/config"
	"github.com/pocka/legit/git"
)

// repositorySummary contains overview of a git repository.
type repositorySummary struct {
	// DisplayName is a directory name without ".git" suffix.
	DisplayName string

	// DirName is a directory name of the repository.
	DirName string

	// Description is a contents of "description" text file in the repository root.
	Description string

	// LastCommitAtRelative is a relative datetime string of last commit.
	// For example, "1 hour ago" or "2 years ago".
	LastCommitAtRelative string

	LastCommit *object.Commit
}

// repoListData is a data object passed to "repo-list" template.
type repoListData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Repositories is a slice of every repositories legit sees.
	Repositories []repositorySummary
}

// repositoryMeta is a shared data object passed to every pages under each repositories.
type repositoryMeta struct {
	// DisplayName is a directory name without ".git" suffix.
	DisplayName string

	// DirName is a directory name of the repository.
	DirName string

	// Description is a contents of "description" text file in the repository root.
	Description string

	// Ref is a ref for the current context. If a page is not tied to refs, default branch
	// will be set.
	Ref string
}

// repoTopData is a data object passed to "repo-top" template.
type repoTopData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Rendered README.
	Readme template.HTML

	// DefaultBranch is textual representation of repository's default branch.
	DefaultBranch string

	// RecentCommits is a list of recent commits made in the default branch.
	RecentCommits []*object.Commit

	// Whether this repository is available as Go Module.
	IsGoModule bool
}

// repoRefsData is a data object passed to "repo-refs" template.
type repoRefsData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Tags is a list of git tags (annotate and lightweight) in the repository.
	Tags []*git.TagReference

	Branches []*plumbing.Reference
}

// repoTreeRefData is a data object passed to "repo-tree-ref" template.
type repoTreeRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Path to the current directory. On repository root, this is empty slice.
	Path []string

	// Files is a list of files for the current directory.
	Files []git.NiceTree
}

// repoBlobRefData is a data object passed to "repo-blob-ref" template.
type repoBlobRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Path to the blob.
	Path []string

	// Content of the blob.
	Content string

	SyntaxHighlightedContent template.HTML

	// LineNumbers holds sequential numbers starting from 1 up to line count of the blob.
	LineNumbers []uint
}

// repoLogRefData is a data object passed to "repo-log-ref" template.
type repoLogRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Commits made to the ref.
	Commits []*object.Commit
}

// repoCommitData is a data object passed to "repo-commit" template.
type repoCommitData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	Commit *object.Commit
	Parent *object.Commit

	Diff *git.NiceDiff
}

// error404Data is a data object passed to "404" template.
type error404Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config
}

// error500Data is a data object passed to "500" template.
type error500Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config
}
