// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

import http from "k6/http";

export const options = {
	iterations: 10,
};

export default function() {
	http.get("http://localhost:8080/foo.git");
}
