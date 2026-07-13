// This package exposes embedded resources across modules.
// A dedicated module is necessary due to Go's restriction on "embed" package
// cannot embed resources on parent directories (and Go not allowing circular
// dependencies.)
//
// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package embed

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var defaultStaticDir embed.FS

//go:embed templates/*
var defaultTemplatesDir embed.FS

// StaticDir returns embedded default static contents directory.
func StaticDir() fs.FS {
	// Go cannot embed directory as-is. Stripping the subdirectory.
	f, err := fs.Sub(defaultStaticDir, "static")
	if err != nil {
		panic("static/ directory not embedded correctly")
	}

	return f
}

// TemplatesDir returns embedded default templates directory.
func TemplatesDir() fs.FS {
	// Go cannot embed directory as-is. Stripping the subdirectory.
	f, err := fs.Sub(defaultTemplatesDir, "templates")
	if err != nil {
		panic("templates/ directory not embedded correctly")
	}

	return f
}
