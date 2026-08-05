//go:build debug

// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT
//
// Debug routes. To access this routes, add a "debug" build tag.

package debug

import (
	"net/http"
	"net/http/pprof"
	"path/filepath"
)

func Handle(w http.ResponseWriter, r *http.Request) bool {
	// We have to use "/debug/pprof" path as the "net/http/pprof" package
	// recklessly hardcodes the path. That problematic package also registers
	// these paths using problematic "init()" calls, but our handlers takes
	// higher precedences, so we have to register manually.
	switch filepath.Clean(r.URL.Path) {
	case "/debug/pprof":
		pprof.Index(w, r)
		return true
	case "/debug/pprof/cmdline":
		pprof.Cmdline(w, r)
		return true
	case "/debug/pprof/profile":
		pprof.Profile(w, r)
		return true
	case "/debug/pprof/symbol":
		pprof.Symbol(w, r)
		return true
	case "/debug/pprof/trace":
		pprof.Trace(w, r)
		return true
	}

	return false
}
