<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Changelog

## v1.0.0

This entry states changes from the original legit.

### New Features

- Completely new default CSS / HTML for better screen real estate utilization.
- Customizable page footer.
- Tab-width selector (requires JavaScript.)
- Pagination for commits (log) page.
- HTML preview for Markdown files on blob page.
- Image display support for Markdown files, both on blob page and summary page.
- New CLI flags: `-server.host`, `-server.port`, `-repo.scanPath`.

### Removed Features

- [go-import](https://go.dev/ref/mod#vcs-find) support: Go cannot handle `.git` extension and legit has no option to trim the extension.
- Syntax highlight theming: HTML no longer contains hardcoded colors, so edit CSS for theme customization.

### Bug Fixes

- Ignored repositories are no longer clonable / fetchable. ([icyphox/legit#56](https://github.com/icyphox/legit/issues/56), [pocka.jp/legit#2](https://tangled.org/pocka.jp/legit/issues/2))
- `git fetch` errors when Git client compress request with Gzip. ([icyphox/legit#58](https://github.com/icyphox/legit/pull/58))
- Fix Nix derivation no declaring `git` dependency.
- Paths inside config file are now resolved from the config file, rather than working directory.

### Performance

- Stop generating HTML sanitizers policy object on each request, slightly improved Markdown and README rendering. ([pocka.jp/legit#4](https://tangled.org/pocka.jp/legit/issues/4))
- Improved the summary page's response time, especially for repositories with many commits, by fixing the code loads an entire commit history for "Recent commits".

### Security

- Use Landlock LSM for filesystem access restriction, similar to unveil(2) on OpenBSD. ([pocka.jp/legit#13](https://tangled.org/pocka.jp/legit/issues/13))

### Other Changes

- Data type of template variables (data passed Go template) are defined in `routes/data.go`.
- You can start legit without `-config` option, as long as you provide `-repo.scanPath` CLI flag.
- Easy to use `Dockerfile` (only web interface. `git clone` and `git fetch` is not available for now.)
- Added NixOS module and Home Manager module to Nix Flake file.
- Default templates and static files are embedded in the binary. You can omit `dirs.templates` and `dirs.static` options. ([icyphox/legit#10](https://github.com/icyphox/legit/issues/10))
- Upgraded chroma package from v2.14 to v2.27. Blob page highlights more languages.
- Upgraded go-git package from v5.6 to v5.13.
