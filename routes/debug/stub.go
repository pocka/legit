//go:build !debug

// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT
//
// This file will be used for non-debug build, to provide no-op handler.
// The purpose of this conditional compilation is to prevent pprof code from
// being included in the non-debug builds. The "net/http/pprof" package uses
// "init" to register default routes carelessly and not importing the file
// is simple and reliable solution.

package debug

import (
	"net/http"
)

func Register(*http.ServeMux) {}
