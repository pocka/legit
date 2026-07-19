{
  description = "web frontend for git";

  inputs.nixpkgs.url = "github:nixos/nixpkgs";

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            inherit system;
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      packages = forAllSystems (
        { system, pkgs }: {
          legit = pkgs.callPackage ./. { };
          default = self.packages.${system}.legit;

          testing = pkgs.callPackage ./nix/test-server.nix { };
        }
      );

      homeManagerModules.default = import ./nix/home-manager-module.nix self;
      nixosModules.default = import ./nix/nixos-module.nix self;

      apps = forAllSystems (
        { system, pkgs }: {
          k6 = {
            type = "app";
            program = pkgs.lib.getExe pkgs.k6;
          };

          dprint = {
            type = "app";
            program = pkgs.lib.getExe pkgs.dprint;
          };
        }
      );

      formatter = forAllSystems (
        { system, pkgs }:
        pkgs.buildFHSEnv {
          name = "fhs-dprint";
          targetPkgs =
            pkgs: with pkgs; [
              # Formatter frontend.
              # https://dprint.dev/
              dprint

              # > Official formatter for Nix code
              # https://github.com/NixOS/nixfmt
              nixfmt

              # For "gofmt" command.
              go
            ];
          runScript = "dprint fmt";
        }
      );

      devShells = forAllSystems (
        { system, pkgs }: {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              go
              gopls
            ];
          };

          debug = pkgs.mkShell {
            # https://github.com/go-delve/delve/issues/3085
            hardeningDisable = [ "fortify" ];
            packages = with pkgs; [
              go
              delve
            ];
          };

          # For those who don't have podman on system.
          podman = pkgs.callPackage ./nix/devshell-podman.nix { };
        }
      );

      nixosConfigurations =
        let
          system = "x86_64-linux";
        in
        {
          soft-legit = nixpkgs.lib.nixosSystem {
            inherit system;
            pkgs = import nixpkgs { inherit system; };

            modules = [
              "${nixpkgs}/nixos/modules/virtualisation/qemu-vm.nix"
              ./nix/nixos-configuration-soft-legit.nix
              self.nixosModules.default
            ];
          };

          soft-legit-hm = nixpkgs.lib.nixosSystem {
            inherit system;
            pkgs = import nixpkgs { inherit system; };

            modules = [
              "${nixpkgs}/nixos/modules/virtualisation/qemu-vm.nix"
              ./nix/nixos-configuration-soft-legit-hm.nix
              ({ ... }: { home-manager.extraSpecialArgs.legit = self; })
            ];
          };
        };
    };
}
