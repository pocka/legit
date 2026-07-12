<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Development

Commands and tips for developing legit.

## Load Testing

`tests/k6` directory contains `*.js` for [k6](https://grafana.com/oss/k6/), an OSS load testing tool.

These test files access a test server on localhost:8080.
Run `nix run .#testing` (you need Nix Flakes) or follow steps described in `tests/k6/.gitignore` manually (you need bash, git, and built `legit` binary) to launch the test server.

Once the test server is ready, pass your desired test file to `k6 run` command.
Nix user can use `nix run .#k6 -- run` without installing k6 manually.

## Code Formatting

This project use [dprint](https://dprint.dev/).
You have to install [`nixfmt`](https://github.com/NixOS/nixfmt) and Go toolchain as well.

Nix user can run `nix fmt` without installing or configuring anything.
