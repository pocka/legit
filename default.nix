# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  buildGoModule,
  lib,
  git,
}:
buildGoModule {
  name = "legit";
  version = "1.0.0";

  src =
    with lib.fileset;
    toSource {
      root = ./.;
      fileset = unions [
        ./go.mod
        ./go.sum
        ./embed
        (fileFilter (file: file.hasExt "go") ./.)
      ];
    };

  vendorHash = "sha256-SWMJVv7QQt4gHaPjb5Q5m20jzFMPHqa+McI26EYg6Ak=";

  ldflags = [
    # git binary from nixpkgs links against libs under "/nix/store/.../lib"
    "-X main.additionalAccessDirs=/nix/store"
    "-X github.com/pocka/legit/git/exe.gitPath=${lib.getExe git}"
  ];

  # Test scripts invoke system "git" command.
  nativeBuildInputs = [ git ];

  meta = {
    mainProgram = "legit";
  };
}
