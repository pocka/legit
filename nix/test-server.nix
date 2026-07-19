# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  callPackage,
  lib,
  git,
  writeShellApplication,
}:
let
  legit = callPackage ../. { };
  repos = callPackage ../tests/k6/repos.nix { };
in
writeShellApplication {
  name = "legit-testing";

  text = ''
    ${lib.getExe legit} -config ${repos}/tests/k6/config.yaml
  '';

  runtimeInputs = [
    git
    legit
    repos
  ];
}
