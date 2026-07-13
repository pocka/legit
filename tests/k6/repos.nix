# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  lib,
  stdenv,
  bash,
  git,
}:
stdenv.mkDerivation {
  name = "legit-test-repos";

  src = lib.fileset.toSource {
    root = ../../.;
    fileset = lib.fileset.unions [
      ./create_repos.bash
      ./config.yaml
    ];
  };

  postInstall = ''
    mkdir -p $out/tests/k6/repos
    cd $out/tests/k6/repos
    bash $src/tests/k6/create_repos.bash
    ln -s $src/tests/k6/config.yaml $out/tests/k6/config.yaml
  '';

  nativeBuildInputs = [
    bash
    git
  ];
}
