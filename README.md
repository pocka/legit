<!--
Copyright 2025 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# legit

legit is a web frontend for git repositories, written in Go.

This code is a fork of <https://github.com/icyphox/legit>, with HTML/CSS/JS customization and typed template pipelines.
I'm using legit along with [soft-serve](https://github.com/charmbracelet/soft-serve) and this fork is optimized for this usecase; good viewing experience, minimum Git operation for security and integration for soft-serve.

## Features

- Repository browsing; commits, refs, tree, etc.
- Simple deployment; single binary without CGI.
- Template customization; modify HTML/CSS/JS to your liking.
- Supports GitWeb `description` file.

## Requirements

- Building and running without building requires Go toolchain >= v1.24.1.
- Both Linux and macOS is supported.
- Put TLS terminating proxy such as nginx or Caddy in front of legit.

## Install

### Manual

Clone and run `go build` at repository root directory.
Go compiler generates `legit` executable file at the repository root directory.

You can also run legit without installing by `go run .`.

### Nix

As a quick and dirty way, you can use the original `legit-web` package and overlay to install this fork.

```nix
final: prev:
{
  legit-web = prev.legit-web.overrideAttrs (old: {
    src = prev.fetchFromGitHub {
      owner = "pocka";
      repo = "legit";
      rev = "bc147a9425e6265adca2672103c0d0b0dfcd735d";
      hash = "sha256-We3ceKWo9viSfM9C/l7CvKiwfGf8bbKvH7M6M0xU1Cg=";
    };

    vendorHash = "sha256-QxkMxO8uzBCC3oMSWjdVsbR2cluYMx5OOKTgaNOLHxc=";
  });
}
```

Runtime error will happen if Go toolchain in your nixpkgs is older than v1.24.1.

## Configuration

legit reads YAML config file. Create YAML file somewhere (e.g. `$XDG_CONFIG_HOME/legit/config.yaml`) and pass the path to legit via `--config` flag.
See the sample [`config.yaml`](./config.yaml) for more info.

## License

This software is licensed under MIT.
See `LICENSE` file for license text.

Newly added files have [REUSE](https://reuse.software/) compliant comment headers for easier per-file use.
