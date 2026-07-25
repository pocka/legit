//go:build debug

// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT
//
// Debug routes. To access this routes, add a "debug" build tag.

package debug

import (
	"net/http"
	"net/http/pprof"
)

func Register(mux *http.ServeMux) {
	// We have to use "/debug/pprof" path as the "net/http/pprof" package
	// recklessly hardcodes the path. That problematic package also registers
	// these paths using problematic "init()" calls, but our handlers takes
	// higher precedences, so we have to register manually.
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
}
