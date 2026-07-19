//go:build !(openbsd || linux)

// Stub functions for GOOS that don't support filesystem restriction.

package main

func restrictFileAccessTo(_dirs ...filesystemAccess) error {
	return nil
}
