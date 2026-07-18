# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

self:
{
  config,
  lib,
  pkgs,
  ...
}:
let
  yaml = pkgs.formats.yaml { };
in
{
  options = {
    services.legit = {
      enable = lib.mkEnableOption "legit";
      package = lib.mkOption {
        type = lib.types.package;
        default = self.packages.${pkgs.stdenv.system}.legit;
      };
      config = lib.mkOption {
        type = yaml.type;
        default = {
          repo.scanPath = "${config.xdg.dataHome}/legit/repos";
          meta = {
            syntaxHighlight = true;
          };
          footer = {
            poweredBy = true;
          };
          server = {
            host = "127.0.0.1";
            port = 5555;
          };
        };
        description = ''
          The contents of the configuration file for legit.
        '';

        example = lib.literalExpression ''
          {
            repo.scanPath = "''${config.xdg.dataHome}/legit/repos";
            meta = {
              title = "My Projects";
              description = "My git repos";
              syntaxHighlight = true;
            };
            footer = {
              links = [
                { text = "Contact"; href = "mailto:alice@example.com"; }
              ];
              poweredBy = true;
            };
            server = {
              name = "git.example.com";
              host = "127.0.0.1";
              port = 5555;
            };
          }
        '';
      };
    };
  };

  config =
    let
      cfg = config.services.legit;
      configFile = yaml.generate "legit.config.yaml" cfg.config;
    in
    lib.mkIf cfg.enable {
      systemd.user.paths.legit-repos = {
        Unit = {
          Description = "repositories directory read by legit";
        };

        Install = {
          WantedBy = [ "default.target" ];
        };

        Path = {
          PathExists = cfg.config.repo.scanPath;
          Unit = "legit.service";
        };
      };

      systemd.user.services.legit = {
        Unit = {
          Description = "web frontend for git repositories";
          After = [ "network.target" ];
        };

        Service = {
          Type = "simple";
          ExecStart = "${lib.getExe cfg.package} -config ${configFile}";
          Restart = "always";
        };
      };
    };
}
