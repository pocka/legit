//go:build openbsd

package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func (a filesystemAccess) mode() string {
	m := ""

	if a.read {
		m = m + "r"
	}

	if a.write {
		m = m + "w"
	}

	if a.execute {
		m = m + "x"
	}

	return m
}

func restrictFileAccessTo(allowList ...filesystemAccess) error {
	for _, a := range allowList {
		if err := unix.Unveil(a.path, a.mode()); err != nil {
			return fmt.Errorf("Unveil error (%s:%s): %w", a.path, a.mode(), err)
		}
	}

	if err := unix.UnveilBlock(); err != nil {
		return fmt.Errorf("Unveil block error: %w", err)
	}

	return nil
}
