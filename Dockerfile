# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

FROM golang:1.26-alpine3.24 AS build
WORKDIR /app

# Download external modules for caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy source files needed for building the binary.
COPY *.go ./
COPY config/*.go ./config/
COPY renderer ./renderer
COPY git ./git
COPY routes ./routes
COPY embed ./embed
RUN go build -o /bin/legit -ldflags="-X 'main.additionalAccessDirs=/usr/lib,/lib'"

# Scratch image does not have mkdir.
RUN mkdir -p /var/www/legit

# ---

FROM alpine/git:2.54.0
WORKDIR /etc/legit

COPY --from=build /bin/legit /bin/legit

# Mount "config.yaml" file at "/etc/legit/config.yaml".
# You can optionally mount "static" and "templates" directories under "/etc/legit/"
# then tell legit to read them.
#
#     dirs:
#       templates: /etc/legit/templates
#       static: /etc/legit/static
COPY config/base.yaml /etc/legit/config.yaml
VOLUME ["/etc/legit"]

# "/var/www/legit" holds git repositories to host. You have to mount at this
# otherwise legit hosts empty top page, which is useless.
COPY --from=build /var/www/legit /var/www/legit
VOLUME ["/var/www/legit"]

# legit serves HTTP content on this TCP port. Bind it to host's :80 or whatever
# you want.
EXPOSE 5555

# Override default entrypoint set by alpine/git.
ENTRYPOINT []

# Prevent git from accessing config files.
ENV GIT_CONFIG_GLOBAL=/dev/null
ENV GIT_CONFIG_SYSTEM=/dev/null

# File paths and listen address depend on how the container was built.
# Users are not supposed to change these options.
CMD [ \
	"/bin/legit", \
	"-config", "/etc/legit/config.yaml", \
	"-server.host", "0.0.0.0", \
	"-server.port", "5555", \
	"-repo.scanPath", "/var/www/legit" \
]
