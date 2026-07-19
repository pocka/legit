//go:build linux

// Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"

	"github.com/landlock-lsm/go-landlock/landlock"
)

func restrictFileAccessTo(allowList ...filesystemAccess) error {
	rules := make([]landlock.Rule, 0, len(allowList))
	for _, a := range allowList {
		var rule landlock.Rule

		if a.isDir {
			if a.write {
				rule = landlock.RWDirs(a.path)
			} else {
				rule = landlock.RODirs(a.path)
			}
		} else {
			if a.write {
				rule = landlock.RWFiles(a.path)
			} else {
				rule = landlock.ROFiles(a.path)
			}
		}

		rules = append(rules, rule)
	}

	err := landlock.V9.BestEffort().RestrictPaths(rules...)

	if err != nil {
		return fmt.Errorf("Landlock error: %w", err)
	}

	return nil
}
