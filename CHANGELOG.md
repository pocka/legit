<!--
Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
SPDX-License-Identifier: MIT
-->

# Changelog

## v1.3.0 - 2026-08-21

### New Features

- Added `ui.email` config option as a countermeasure against email spammers.
- Raw page for SVG files now returns `image/svg+xml` rather than `text/plain`, to allow browsers to render the image.

### Bug Fixes

- Fixed certain parent commits are not displayed for a multiple parent commit. (e.g., commit created by `jj new x y`)
- Fixed relative link in HTML preview of Markdown files inside a directory (e.g., `foo/bar.md`) pointing a nonexistent path.

### Security

- Return `X-Content-Type-Options: nosiff` HTTP header on raw pages to prevent browsers from MIME sniffing (change content type based on file content, instead of the `Content-Type` returned by server.)

### Other Changes

- Opening a link to another Markdown file now opens HTML preview of the file, rather than source code view.
- Email address of commit author and committer is not available on HTML. Let client-side JavaScript restores the data, or you can restore the previous spammer-friendly behavior by setting `ui.email` to `"mailto"`.

## v1.2.1 - 2026-08-16

### Security

- Fixed XSS in raw preview of attacker controlled HTML file.

## v1.2.0 - 2026-08-08

### New Features

- Added `repo.diff` config option / CLI flag to select `go` (pure Go diff implementation, default) or `system` (system install git, much faster) for generating diffs.
- Added `ui.diff.hideThresholdLines` config option to collapse large diffs.
- Added `repo.trimDotGitSuffix` config option to serve bare repositories without `.git` suffix. You can still access `/<name>.git/*` paths when this option is enabled.
- Added `meta.robotsTxt` config option to specify a file for `/robots.txt` requests.

### Bug Fixes

- Fixed pages for a ref including slash character (e.g., `feature/foo`) not working.
- Removed double-slash from links on tree page.
- Fixed header section overlaps hash-scroll target (e.g., latest commit link on a log page.)

### Other Changes

- Reduced reflow / screen flickering on commit page.

## v1.1.0 - 2026-08-02

### New Features

- Grouping repositories by gitweb-compatible category (opt-in, see the added `ui.category.*` options). For how to set a category for a repository, wefer to [gitweb's document](https://git-scm.com/docs/gitweb#Documentation/gitweb.txt-categoryorgitwebcategory).
- Added "Preferences" dialog to let users (visitors) allow customize theme color and file list appearance.
- Tooltip on commit / author datetime now has relative time. E.g., 1 year ago
- When `$DIR_DIR/description` file is absent, legit now reads `gitweb.description` git config and use it as a description text, for gitweb compatibility.
- Added `staticDirRevision` config / build flag for browser-cache busting.

### Bug Fixes

- Fixed `git clone` and `git fetch` are unavailable on OCI image.
- Fixed ignored repository reveal using URL-encoded slash character, such as `.%2Fprivate_repo`.
- Fixed `description` file in a checked out worktree being treated as a description file.
- Fixed broken layout of 404/500 pages.

### Security

- Switched to more robust repository directory path check for more secure filesystem access outside OpenBSD and Linux.

### Other Changes

- Tab width selector has been moved to the new "Preferences" dialog.
- Made the repository list page (`/`) cache-friendly by outputting timestamp rather than relative time. If user enabled JS in a browser, this change is effectively no-op.
- Added `-compileTemplatesOnRequest`, `-dirs.static` and `-dirs.templates` CLI flags.
- Fixed links that always 307 redirects due to lack of trailing slash.

## v1.0.0 - 2026-07-19

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
