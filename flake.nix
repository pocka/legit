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
        { system, pkgs }:
        let
          legit = self.packages.${system}.legit;
          files = pkgs.lib.fileset.toSource {
            root = ./.;
            fileset = pkgs.lib.fileset.unions [
              ./config.yaml
              ./static
              ./templates
            ];
          };
        in
        {
          legit = pkgs.buildGoModule {
            name = "legit";
            rev = "master";
            src = ./.;

            vendorHash = "sha256-GR+ddG6058s9NDGwbk2z86vcOIa11DUQqP6tS0ruUdo=";
          };

          default = legit;

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "sini:5000/legit";
            tag = "latest";
            contents = [
              files
              legit
              pkgs.git
            ];
            config = {
              Entrypoint = [ "${legit}/bin/legit" ];
              ExposedPorts = {
                "5555/tcp" = { };
              };
            };
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
        }
      );
    };
}
