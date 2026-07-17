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
            fileset = pkgs.lib.fileset.unions [ ./config.yaml ];
          };
        in
        {
          legit = pkgs.buildGoModule {
            name = "legit";
            rev = "master";
            src =
              with pkgs.lib.fileset;
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

            nativeBuildInputs = with pkgs; [ git ];

            meta.mainProgram = "legit";
          };

          default = legit;

          testing =
            let
              repos = pkgs.callPackage ./tests/k6/repos.nix { };
            in
            pkgs.writeShellApplication {
              name = "legit-testing";

              text = ''
                ${pkgs.lib.getExe legit} -config ${repos}/tests/k6/config.yaml
              '';

              runtimeInputs = with pkgs; [
                git
                legit

                repos
              ];
            };
        }
      );

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

          # For those who don't have podman on system.
          podman =
            let
              conf = pkgs.stdenv.mkDerivation {
                name = "podman-config-files";

                # No source.
                unpackPhase = "true";

                postInstall = ''
                  mkdir $out
                  cat > $out/registries.conf <<EOF
                  unqualified-search-registries = ["docker.io"]
                  EOF

                  # This is meant to be copied under ~/.config/containers/
                  cat > $out/policy.json <<EOF
                  {
                    "default": [{ "type": "insecureAcceptAnything" }]
                  }
                  EOF

                  cat > $out/containers.conf <<EOF
                  [engine]
                  cgroup_manager="cgroupfs"
                  events_logger="file"
                  EOF
                '';
              };
            in
            pkgs.mkShell {
              packages = with pkgs; [
                podman
                conf
              ];

              shellHook = ''
                if [[ ! -f ~/.config/containers/policy.json && ! -f /etc/containers/policy.json ]]; then
                  echo "Create policy.json file for podman."
                  echo "https://podman.io/docs/installation#policyjson"
                  echo "If you are okay with insecureAcceptAnything for all, run:"
                  echo "install -Dm555 ${conf}/policy.json ~/.config/containers/policy.json"
                fi
              '';

              CONTAINERS_REGISTRIES_CONF = "${conf}/registries.conf";
              CONTAINERS_CONF = "${conf}/containers.conf";
            };
        }
      );
    };
}
