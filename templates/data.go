// This package defines the data passed to templates.
//
// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package templates

import (
	"fmt"
	"html/template"
	"maps"
	"slices"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/git"
)

// RepositorySummary contains overview of a git repository.
type RepositorySummary struct {
	// DisplayName is a directory name without ".git" suffix.
	DisplayName string

	// DirName is a directory name of the repository.
	DirName string

	// Description is a contents of "description" text file in the repository root.
	Description string

	// Category is gitweb compatible repository category text. Unless "ui.category.grouping"
	// option is set, this field is always empty.
	Category string

	LastCommit *object.Commit
}

// RepoListData is a data object passed to "repo-list" template.
type RepoListData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Repositories is a slice of every repositories legit sees.
	Repositories []RepositorySummary
}

type RepositoriesByCategory struct {
	Category     string
	Repositories []RepositorySummary
}

func (d RepoListData) RepositoriesByCategory() []RepositoriesByCategory {
	if !d.Config.UI.Category.Grouping {
		return []RepositoriesByCategory{
			{
				Category:     "",
				Repositories: d.Repositories,
			},
		}
	}

	m := make(map[string]*RepositoriesByCategory, len(d.Config.UI.Category.Order))

	for _, repo := range d.Repositories {
		category := repo.Category
		if category == "" {
			category = d.Config.UI.Category.Default
		}

		if group, ok := m[category]; !ok {
			m[category] = &RepositoriesByCategory{
				Category:     category,
				Repositories: []RepositorySummary{repo},
			}
		} else {
			group.Repositories = append(group.Repositories, repo)
		}
	}

	out := make([]RepositoriesByCategory, 0, len(m))

	if empty, ok := m[""]; ok {
		out = append(out, *empty)
	}

	for _, prioritized := range d.Config.UI.Category.Order {
		if group, ok := m[prioritized]; ok {
			out = append(out, *group)
		}
	}

	for _, category := range slices.Sorted(maps.Keys(m)) {
		if category == "" || slices.Contains(d.Config.UI.Category.Order, category) {
			continue
		}

		out = append(out, *m[category])
	}

	return out
}

// RepositoryMeta is a shared data object passed to every pages under each repositories.
type RepositoryMeta struct {
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

// RepoTopData is a data object passed to "repo-top" template.
type RepoTopData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Rendered README.
	Readme template.HTML

	// DefaultBranch is textual representation of repository's default branch.
	DefaultBranch string

	// RecentCommits is a list of recent commits made in the default branch.
	RecentCommits []*object.Commit

	// Whether this repository is available as Go Module.
	IsGoModule bool
}

// RepoRefsData is a data object passed to "repo-refs" template.
type RepoRefsData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Tags is a list of git tags (annotate and lightweight) in the repository.
	Tags []*git.TagReference

	Branches []*plumbing.Reference
}

// RepoTreeRefData is a data object passed to "repo-tree-ref" template.
type RepoTreeRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Path to the current directory. On repository root, this is empty slice.
	Path []string

	// Files is a list of files for the current directory.
	Files []git.NiceTree
}

// RepoBlobRefData is a data object passed to "repo-blob-ref" template.
type RepoBlobRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Path to the blob.
	Path []string

	// Content of the blob.
	Content string

	SyntaxHighlightedContent template.HTML

	// LineNumbers holds sequential numbers starting from 1 up to line count of the blob.
	LineNumbers []uint

	// A list of preview output types.
	PreviewTypes []string
}

// RepoBlobRefHTMLPreviewData is a data object passed to "repo-blob-ref-html-preview" template.
type RepoBlobRefHTMLPreviewData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Path to the blob.
	Path []string

	// Rendered HTML
	Content template.HTML
}

// RepoLogRefData is a data object passed to "repo-log-ref" template.
type RepoLogRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	// Commits made to the ref.
	Commits []*object.Commit

	PrevPageHref string
	NextPageHref string
}

func (r RepoLogRefData) TopCommitForRange() *object.Commit {
	if len(r.Commits) == 0 {
		return nil
	}

	return r.Commits[0]
}

func (r RepoLogRefData) BottomCommitForRange() *object.Commit {
	if len(r.Commits) == 0 {
		return nil
	}

	return r.Commits[len(r.Commits)-1]
}

// RepoCommitData is a data object passed to "repo-commit" template.
type RepoCommitData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta RepositoryMeta

	Commit *object.Commit
	Parent *object.Commit

	Diff *git.NiceDiff
}

func (r RepoCommitData) IsLargeDiff(file *gitdiff.File) bool {
	if r.Config.UI.Diff.HideThresholdLines == 0 {
		return false
	}

	var n uint32

	for _, fragment := range file.TextFragments {
		n += uint32(max(0, fragment.LinesAdded))
		n += uint32(max(0, fragment.LinesDeleted))
	}

	return n > r.Config.UI.Diff.HideThresholdLines
}

func (r RepoCommitData) FileDiffStat(file *gitdiff.File) string {
	var added int64
	var deleted int64

	for _, fragment := range file.TextFragments {
		added += fragment.LinesAdded
		deleted += fragment.LinesDeleted
	}

	return fmt.Sprintf("+%d/-%d", added, deleted)
}

// Error404Data is a data object passed to "404" template.
type Error404Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Diagnosis is debug text supposed to attach to a bug report for website owner.
	Diagnosis string
}

// Error500Data is a data object passed to "500" template.
type Error500Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Diagnosis is debug text supposed to attach to a bug report for website owner.
	Diagnosis string
}
