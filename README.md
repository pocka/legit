<!--
Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# legit

legit is a readonly web frontend for git repositories, written in [Go](https://go.dev/).

This code is a hard fork of <https://github.com/icyphox/legit>.
Various bugs are fixed and templates (HTML/CSS/JS) are completely different.
New features are opt-in, so you can use a config file from the upstream in this fork.

## Features

- Repository browsing; commits, refs, tree, etc.
- Simple deployment; single binary without CGI.
- Supports [gitweb](https://git-scm.com/docs/gitweb)-compatible description and category.
- Secure readonly file access with [unveil(2)](https://man.openbsd.org/unveil.2) (OpenBSD) and [Landlock](https://docs.kernel.org/userspace-api/landlock.html) (Linux).

## Requirements

- Supports Linux and OpenBSD. Probably runs on other UNIX-y systems as well, but lacks important security features.
- Put TLS terminating proxy such as reverse proxy or CDN in front of legit.

## Quick Start

Clone this repository and run `go run .` on the worktree.
You need Go toolchain >= v1.25.0 to run `go` commands.

```sh
go run . -repo.scanPath demo
```

See [docs/INSTALL.md](./docs/INSTALL.md) for more info.

## Configuration

legit reads YAML config file. Create YAML file somewhere (e.g. `$XDG_CONFIG_HOME/legit/config.yaml`) and pass the path to legit via `--config` flag.
See the sample [`config.yaml`](./config.yaml) for more info.

## Bug Reports

If you find a bug in this software, please report it on <https://tangled.org/pocka.jp/legit/issues>.

## License

This software is licensed under MIT.
See [`license`](./license) for license text.

Newly added files have [REUSE](https://reuse.software/) compliant comment headers for easier per-file use.
