<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Install Guide

legit is a single binary [Go](https://go.dev/) application.
Manually build and place the built binary to your liking, or run / install with [Nix](https://nix.dev/manual/nix/2.34/introduction.html) [Flake](https://wiki.nixos.org/wiki/Flakes) or OCI tool ([Docker](https://www.docker.com/), [Podman](https://podman.io/)).

## Manual Build

legit is a Go project.
You need Go toolchain >= v1.25.0 for building the server binary.

If you're building for Linux, get your system's dynamic libraries directory before building the binary.
See "Limiting file access" section on [SECURITY.md](SECURITY.md) for more info.

At the root directory of this repository, run:

```sh
go build
```

to generate an executable file named `legit`.
You're now ready to launch a legit server.

## Nix Flake

The default package of this repository's Flake is a derivation of a legit server.

```sh
# Serves ./repos directory at http://localhost:5555
nix run git+https://git.pocka.jp/legit -- -repo.scanPath ./repos
```

In addition to that, the Flake file exports a [NixOS](https://nixos.org/) module and a [Home Manager](https://github.com/nix-community/home-manager) module.
Both modules accept `config` attrset, which will be simply converted to a YAML file.

`nix/` directory contains sample NixOS configurations [`nixos-configuration-soft-legit.nix`](../nix/nixos-configuration-soft-legit.nix) (NixOS module) and [`nixos-configurations-soft-legit-hm.nix](../nix/nixos-configuration-soft-legit-hm.nix) (Home Manager module).

Both modules use systemd path unit to launch the main service.
legit service won't start until you create the `repo.scanPath` directory.

## OCI

This fork does not publish an OCI image to online registries, so you have to build yourself.

Generated image,

- exposes TCP port 5555 for HTTP server, and
- expects `/var/www/legit` to be a mounted directory that contains git repositories

To configure, mount YAML config file at `/etc/legit/config.yaml`.

Here is example commands for building and running a legit OCI image:

```sh
podman build . -t pocka/legit
podman -v ./repos:/var/www/legit -v ./config.yaml:/etc/legit/config.yaml --publish 5555:5555 pocka/legit
```

The OCI image ignores these config options:

- `server.host`
- `server.port`
- `repo.scanPath`
