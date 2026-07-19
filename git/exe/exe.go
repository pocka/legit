// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package exe

var gitPath string

func GitBin() string {
	if gitPath != "" {
		return gitPath
	}

	// Let the system resolve.
	return "git"
}
