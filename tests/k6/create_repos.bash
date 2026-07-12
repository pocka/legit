# create_repos.bash
# Bash script to create sample repos in the current directory.
#
# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

# Usage
#   cd test_dir
#   bash /path/to/create_repos.bash

set -e

# To get deterministic output.
export GIT_CONFIG_GLOBAL=/dev/null
export GIT_CONFIG_SYSTEM=/dev/null
export GIT_COMMITTER_DATE="2026-06-30T09:00:00+09:00"

GIT_FLAGS=(
	"-c" "user.name=Alice"
	"-c" "user.email=alice@example.com"
	"-c" "init.defaultBranch=trunk"
)

# foo - basic sample
git ${GIT_FLAGS[@]} init foo

cat <<EOF > foo/README.md
# Foo

This is foo.
I'm [Markdown](https://commonmark.org/) file for *load* **testing**.

```
This file is to test Markdown rendering performance of summary page.
```

* longer
* the
* contents
* and
* more
  + readable
	+ comparable
* output
EOF
git ${GIT_FLAGS[@]} -C foo add README.md
git ${GIT_FLAGS[@]} -C foo commit \
	-m "Add README" \
	--date "2026-06-06T12:00:00+09:00"

cat <<EOF > foo/main.go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}
EOF
git ${GIT_FLAGS[@]} -C foo add main.go
git ${GIT_FLAGS[@]} -C foo commit \
	-m "Create Go application entrypoint" \
	--date "2026-06-06T12:01:00+09:00"

git clone --bare foo foo.git
rm -rf foo
