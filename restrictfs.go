// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package main

type filesystemAccess struct {
	path  string
	isDir bool

	// read allows read operations.
	read bool

	// write allows write operations.
	write bool

	// execute allows execute operations (run a process).
	execute bool
}
