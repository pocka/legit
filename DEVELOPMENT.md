<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Development

Commands and tips for developing legit.

## Issues / Tickets

See <https://tangled.org/pocka.jp/legit/issues> for active issues.

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

## Build an OCI Image

You can build an [OCI](https://opencontainers.org/) image easily with [podman](https://podman.io/).
legit works perfectly fine in an unprivilege (rootless) container.

On the project root directory run the following command:

```sh
podman build . -t pocka/legit
```

That would create and register `pocka/legit` image on your local registry.
To test the image works, run the following command:

```sh
podman run --volume ./demo:/var/www/legit --publish 5555:5555 pocka/legit:latest
```

If you don't have podman on your NixOS system and quickly test these steps, use `nix develop .#podman` devShell.
It has podman and configures minimum podman environment.

## Testing NixOS module

`flake.nix` has [soft-serve](https://github.com/charmbracelet/soft-serve) + [Caddy](https://caddyserver.com/) + legit NixOS configuration for testing purpose.

Run, `nix run .#nixosConfigurations.soft-legit.config.system.build.vm` to launch it in QEMU.
It'll create `nixos.qcow2` file as a disk image, so delete the file to start from scratch.
To access the VM via SSH and HTTP, set `QEMU_NET_OPTS="hostfwd=tcp::<HOST HTTP PORT>-:80,hostfwd=tcp::<HOST SSH PORT>-:22222"` to forward TCP packets.
For example,

```sh
QEMU_NET_OPTS="hostfwd=tcp::8000-:80,hostfwd=tcp::22222-:22222" nix run .#nixosConfigurations.soft-legit.config.system.build.vm
```

As soft-serve has no way to setup initial users without setting a `SOFT_SERVE_INITIAL_ADMIN_KEYS` environment variable _on first launch_, you have to do some manual work for initial launch.
The systemd service for soft-service does not start until SSH public keys for initial admin are registered.

```sh
# inside VM...

# Replace https://codeberg.org/pocka.keys with another provider / your account.
curl -sL https://codeberg.org/pocka.keys | write-initial-admin-keys

# ...the "write-initial-admin-keys" script writes to a systemd environment file.
```

Once you write your SSH public key to the file, the service automatically launches and you can access soft-serve as an admin user.
Then, import a git repository via SSH:

```sh
# on host...

# This command uses the host port number I used in the above QEMU launch example command.
ssh -o IdentitiesOnly=yes -p 22222 localhost -- repo import legit https://github.com/pocka/legit
```

After successful repository creation, legit service automatically starts and you can access legit pages on localhost at the forwarded HTTP port.

You can also check NixOS configurations without running the VM, by `nix flake check`.

## Testing Home Manager module

You can test [Home Manager](https://github.com/nix-community/home-manager) module using NixOS VM as well.

See the above "Testing NixOS module" section for required steps.
Key differences are:

- The flake reference will be `.#nixosConfigurations.soft-legit-hm` (`nix run .#nixosConfigurations.soft-legit-hm.system.build.vm`)
- Login user will be `alice` rather than `admin`
- Both soft-serve and legit run under `alice` (you have to add `--user` flag to systemd, like `systemd --user status legit.service`)
