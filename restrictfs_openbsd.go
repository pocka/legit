//go:build openbsd

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func restrictFileAccessTo(dirs ...string) error {
	for _, dir := range dirs {
		if err := unix.Unveil(dir, "r"); err != nil {
			return fmt.Errorf("Unveil error (%s): %w", dir, err)
		}
	}

	if err := unix.UnveilBlock(); err != nil {
		return fmt.Errorf("Unveil block error: %w", err)
	}

	return nil
}
