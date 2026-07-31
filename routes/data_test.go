// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package routes

import (
	"slices"
	"testing"

	"github.com/pocka/legit/config"
)

func TestRepositoriesByCategorySkipsByDefault(t *testing.T) {
	cfg := config.Config{}
	repos := []repositorySummary{
		{
			DisplayName: "Foo",
			Category:    "B",
		},
		{
			DisplayName: "Bar",
			Category:    "A",
		},
		{
			DisplayName: "Baz",
			Category:    "A",
		},
		{
			DisplayName: "Qux",
			Category:    "",
		},
	}

	data := repoListData{
		Config:       &cfg,
		Repositories: repos,
	}

	got := data.RepositoriesByCategory()

	if len(got) != 1 {
		t.Fatalf("Expected single category, got %d", len(got))
	}

	if len(got[0].Repositories) != len(repos) {
		t.Fatalf("Expected %d repositories, got %d", len(repos), len(got[0].Repositories))
	}

	if !slices.Equal(repos, got[0].Repositories) {
		t.Error("Returned repositories are mutated")
	}
}

func TestRepositoriesByCategoryOK(t *testing.T) {
	cfg := config.Config{}
	cfg.UI.Category.Grouping = true
	cfg.UI.Category.Order = []string{"A"}

	repos := []repositorySummary{
		{
			DisplayName: "Foo",
			Category:    "B",
		},
		{
			DisplayName: "Bar",
			Category:    "A",
		},
		{
			DisplayName: "Baz",
			Category:    "A",
		},
		{
			DisplayName: "Qux",
			Category:    "",
		},
	}

	data := repoListData{
		Config:       &cfg,
		Repositories: repos,
	}

	got := data.RepositoriesByCategory()

	if len(got) != 3 {
		t.Fatalf("Expected three categories, got %d", len(got))
	}

	expected := []repositoriesByCategory{
		{
			Category: "",
			Repositories: []repositorySummary{
				{DisplayName: "Qux", Category: ""},
			},
		},
		{
			Category: "A",
			Repositories: []repositorySummary{
				{DisplayName: "Bar", Category: "A"},
				{DisplayName: "Baz", Category: "A"},
			},
		},
		{
			Category: "B",
			Repositories: []repositorySummary{
				{DisplayName: "Foo", Category: "B"},
			},
		},
	}

	if !slices.EqualFunc(got, expected, func(e1, e2 repositoriesByCategory) bool {
		return e1.Category == e2.Category && slices.Equal(e1.Repositories, e2.Repositories)
	}) {
		t.Error("Unexpected output")
	}
}
