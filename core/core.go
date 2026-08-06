// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package core

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"github.com/pocka/legit/config"
	"github.com/pocka/legit/embed"
	"github.com/pocka/legit/renderer/html"
)

const (
	dotGitSuffix = ".git"
)

var (
	ErrRepositoryIsIgnored = errors.New("repository is ignored")
)

type Core struct {
	// Config used to instantiate this core.
	Config *config.Config

	StaticDir fs.FS

	ScanDir *os.Root

	Markdown  html.MarkdownRenderer
	Plaintext html.PlaintextRenderer

	// t stores compiled templates for caching.
	tmpl         *template.Template
	templatesDir fs.FS
}

func New(config *config.Config) (*Core, error) {
	var err error

	ugcPolicy := bluemonday.UGCPolicy()

	core := Core{
		Config:    config,
		Markdown:  html.NewMarkdownRenderer(ugcPolicy),
		Plaintext: html.NewPlaintextRenderer(ugcPolicy),
	}

	if config.Dirs.Static != "" {
		root, err := os.OpenRoot(config.Dirs.Static)
		if err != nil {
			return nil, fmt.Errorf("unable to open static dir: %w", err)
		}

		core.StaticDir = root.FS()
	} else {
		core.StaticDir = embed.StaticDir()
	}

	if config.Dirs.Templates != "" {
		core.templatesDir = os.DirFS(config.Dirs.Templates)
	} else {
		core.templatesDir = embed.TemplatesDir()
	}

	if !config.CompileTemplatesOnRequest {
		core.tmpl = template.Must(template.ParseFS(core.templatesDir, "*"))
	}

	core.ScanDir, err = os.OpenRoot(config.Repo.ScanPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open scan dir: %w", err)
	}

	return &core, nil
}

// Template returns compiled HTML templates.
func (core *Core) Template() *template.Template {
	if core.tmpl != nil {
		return core.tmpl
	}

	return template.Must(template.ParseFS(core.templatesDir, "*"))
}

// RepositoryPath finds a repository by name and returns its filepath.
// name is not supposed to contain slash characters.
func (core *Core) RepositoryPath(name string) (string, error) {
	stat, err := core.ScanDir.Stat(name)
	if err != nil {
		if !core.Config.Repo.TrimDotGitSuffix || strings.HasSuffix(name, dotGitSuffix) {
			return "", err
		}

		stat, err = core.ScanDir.Stat(name + dotGitSuffix)
		if err != nil {
			return "", err
		}
	}

	dirname := stat.Name()
	if slices.Contains(core.Config.Repo.Ignore, dirname) {
		return "", ErrRepositoryIsIgnored
	}

	return filepath.Join(core.ScanDir.Name(), dirname), nil
}

// RepositoryName returns formatted display name of a repository.
func (_ *Core) RepositoryName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), ".git")
}
