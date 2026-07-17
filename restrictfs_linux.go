//go:build linux

// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/landlock-lsm/go-landlock/landlock"
)

func restrictFileAccessTo(dirs ...string) error {
	err := landlock.V9.BestEffort().RestrictPaths(landlock.RODirs(dirs...))

	if err != nil {
		return fmt.Errorf("Landlock error: %w", err)
	}

	return nil
}
