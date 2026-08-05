// This file defines the data passed to templates.
//
// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

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

// repositorySummary contains overview of a git repository.
type repositorySummary struct {
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

// repoListData is a data object passed to "repo-list" template.
type repoListData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Repositories is a slice of every repositories legit sees.
	Repositories []repositorySummary
}

type repositoriesByCategory struct {
	Category     string
	Repositories []repositorySummary
}

func (d repoListData) RepositoriesByCategory() []repositoriesByCategory {
	if !d.Config.UI.Category.Grouping {
		return []repositoriesByCategory{
			{
				Category:     "",
				Repositories: d.Repositories,
			},
		}
	}

	m := make(map[string]*repositoriesByCategory, len(d.Config.UI.Category.Order))

	for _, repo := range d.Repositories {
		category := repo.Category
		if category == "" {
			category = d.Config.UI.Category.Default
		}

		if group, ok := m[category]; !ok {
			m[category] = &repositoriesByCategory{
				Category:     category,
				Repositories: []repositorySummary{repo},
			}
		} else {
			group.Repositories = append(group.Repositories, repo)
		}
	}

	out := make([]repositoriesByCategory, 0, len(m))

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

	// A list of preview output types.
	PreviewTypes []string
}

// repoBlobRefHTMLPreviewData is a data object passed to "repo-blob-ref-html-preview" template.
type repoBlobRefHTMLPreviewData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Path to the blob.
	Path []string

	// Rendered HTML
	Content template.HTML
}

// repoLogRefData is a data object passed to "repo-log-ref" template.
type repoLogRefData struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	Meta repositoryMeta

	// Commits made to the ref.
	Commits []*object.Commit

	PrevPageHref string
	NextPageHref string
}

func (r repoLogRefData) TopCommitForRange() *object.Commit {
	if len(r.Commits) == 0 {
		return nil
	}

	return r.Commits[0]
}

func (r repoLogRefData) BottomCommitForRange() *object.Commit {
	if len(r.Commits) == 0 {
		return nil
	}

	return r.Commits[len(r.Commits)-1]
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

func (r repoCommitData) IsLargeDiff(file *gitdiff.File) bool {
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

func (r repoCommitData) FileDiffStat(file *gitdiff.File) string {
	var added int64
	var deleted int64

	for _, fragment := range file.TextFragments {
		added += fragment.LinesAdded
		deleted += fragment.LinesDeleted
	}

	return fmt.Sprintf("+%d/-%d", added, deleted)
}

// error404Data is a data object passed to "404" template.
type error404Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Diagnosis is debug text supposed to attach to a bug report for website owner.
	Diagnosis string
}

// error500Data is a data object passed to "500" template.
type error500Data struct {
	// Config represents a resolved config based on "config.yaml".
	Config *config.Config

	// Diagnosis is debug text supposed to attach to a bug report for website owner.
	Diagnosis string
}
