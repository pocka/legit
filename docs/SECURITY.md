<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Securing Your Instance

This document describes legit's security features and deployment advice (non-features.)

## Limiting file access

legit uses [unveil(2)](https://man.openbsd.org/unveil.2) on OpenBSD and [Landlock](https://docs.kernel.org/userspace-api/landlock.html) on Linux to limit its own file accesses.

After loading and validating a config file, legit tells OS what filepaths it's going to read.
OS rejects filesystem access from legit such as,

- reads outside the declared files / directories, for example, `/var/log`
- writes to any file (except `/dev/null`, due to Go's implementation quirks)
- exec / fork of non-declared files (OpenBSD only)

If you find unexpected file access error for directory or file that legit _should_ have access to, add the directory or file to `main.additionalAccessDirs` ldflag.
For example, the below build command allows readonly access to `/opt/foo` directory and `/opt/bar` directory:

```sh
go build -ldflags "-X 'main.additionalAccessDirs=/opt/foo,/opt/bar'"
```

If a directory you specified does not exist on startup, legit exits with error status code.
The flag accepts comma-separated list of directories.

### Allow invocation of git binary on Linux

**Only for Linux manual builds.**
If you use Nix Flake or OCI image, or run legit on OpenBSD, you can skip this section.

By default, legit allows read access to `git` command.
However, unless your system has statically linked `git` command, Landlock rejects invocation of the `git` command due to access to dynamically linked libraries.
As it's virtually impossible to pre-specify paths to dynamic libraries because of differences in Linux distributions, you're responsible for finding and setting `main.additionalAccessDirs`.

One way to figure out which directory `git` command loads is to use `ldd`:

```sh
ldd $(which git)
```

The above command outputs dynamically linked libraries and their paths.
Then, set the library path to `main.additionalAccessDirs` and build:

```
go build -ldflags="-X 'main.additionalAccessDirs=/usr/lib,/lib'"
```

You can narrow the directory list further, depends on your threat model.

## Limit network requests

legit does not have an ability to throttle or limit requests.
You should configure these functionality in downstream, such as reverse proxy and CDN.

How to configure these features are outside of this project's scope.
You're responsible for measuring, assessing, deciding, and configuring network-level security.

## Assigning service user

As a general advice, you should run legit as a restricted service user.
To help you configure a user, here is the list of what legit will do:

- Read files and directories supplied via config files and/or CLI flags.
- Read and write to `/dev/null`.
- Exec `git` executable from `$PATH` (or from `github.com/pocka/legit/git/exe.gitPath` if set).
- Start a TCP server at the specified host and port.
- Writes to stdout and stderr.

## Hiding private repositories

Although this fork still supports `repo.ignore` config option for compatibility, you should use OS filesystem and permissions mechanism to hide private repositories.

### Directory convention

The easiest way to hide private repositories is not putting them in the directory at all.
If you're using [soft-serve](https://github.com/charmbracelet/soft-serve), this is the best method as soft-serve does not use system users.

You create a directory for public sharing (e.g., `public/`, `x/`) then place only the repositories you want to share.

```
my-repositories
├── public
│   ├── foo.git
│   └── bar.git
└── baz.git
```

If you're using soft-serve, append the directory prefix using `repo rename` command.

### File permissions

If you have git SSH access using system user, using user / group for visibility control eliminates additional management overhead.

Create or assign a user for legit service (daemon), then configure repository directories' file permission.
