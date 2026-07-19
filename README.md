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

legit uses unveil(2) on OpenBSD and Landlock LSM on Linux.
If you find unexpected filesystem permission error, add the _directory_ to `main.additionalAccessDirs` ldflag.
That flag takes comma-separated list of directories, and unveil/Landlock allows a readonly access to that paths.

### Nix

Add this repository as a Flake input and use `nixosModules.default` or `homeManagerModules.default`.

```nix
# your/flake.nix
{
	inputs = {
		# --- snip ---

		legit = {
			url = "github:pocka/legit";
			inputs.nixpkgs.follows = "nixpkgs"; # optional
		};
	};

	outputs =
		{ nixpkgs, legit, ... }:
		{
			nixosConfigurations.foo = nixpkgs.lib.nixosSystem {
				# --- snip ---

				modules = [
					# --- snip ---
					legit.nixosModules.default
				];
			};
		};
}
```

### OCI (Docker, Podman)

Build the image on the project root.
The image exposes TCP port 5555 for HTTP server.

Generated image expectes two volume mounts:

- `/var/www/legit` ... a directory containing git repositories to host.
- `/etc/legit/config.yaml` ... config file for UI and metadata customization.

Example commands using podman:

```sh
podman build . -t pocka/legit
podman run -v ./demo:/var/www/legit -v ./config.yaml:/etc/legit/config.yaml --publish 5555:5555 pocka/legit:latest
```

OCI image entrypoint overwrites these config options, so values inside your `config.yaml` will be ignored:

- `server.host`
- `server.port`
- `repo.scanPath`

## Configuration

legit reads YAML config file. Create YAML file somewhere (e.g. `$XDG_CONFIG_HOME/legit/config.yaml`) and pass the path to legit via `--config` flag.
See the sample [`config.yaml`](./config.yaml) for more info.

## License

This software is licensed under MIT.
See `LICENSE` file for license text.

Newly added files have [REUSE](https://reuse.software/) compliant comment headers for easier per-file use.
