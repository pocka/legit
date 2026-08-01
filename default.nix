# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  buildGoModule,
  lib,
  git,
}:
buildGoModule (finalAttrs: {
  name = "legit";
  version = "1.1.0";

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

  vendorHash = "sha256-mSr8uddh7J9P0BhYH7D6riW3KGmr5Qf8caMhKwoUCO0=";

  ldflags = [
    # git binary from nixpkgs links against libs under "/nix/store/.../lib"
    "-X main.additionalAccessDirs=/nix/store"
    "-X github.com/pocka/legit/git/exe.gitPath=${lib.getExe git}"
    "-X github.com/pocka/legit/config.staticDirRevision=v${finalAttrs.version}-embed"
  ];

  # Test scripts invoke system "git" command.
  nativeBuildInputs = [ git ];

  meta = {
    mainProgram = "legit";
  };
})
