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
  disabledModules = [ "services/networking/legit.nix" ];

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
          repo.scanPath = "/var/www/git";
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
            repo.scanPath = "/var/www/git";
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

      user = lib.mkOption {
        type = lib.types.nonEmptyStr;
        default = "legit";
        description = "User account under which legit runs.";
      };

      group = lib.mkOption {
        type = lib.types.nonEmptyStr;
        default = "legit";
        description = "Group account under which legit runs.";
      };
    };
  };

  config =
    let
      cfg = config.services.legit;

      configFile = yaml.generate "legit.config.yaml" cfg.config;
    in
    lib.mkIf cfg.enable {
      users.groups = lib.optionalAttrs (cfg.group == "legit") { legit = { }; };

      users.users = lib.optionalAttrs (cfg.group == "legit") {
        legit = {
          group = "legit";
          isSystemUser = true;
        };
      };

      systemd.services.legit = {
        description = "web frontend for git repositories";
        after = [ "network.target" ];
        restartTriggers = [ configFile ];

        serviceConfig = {
          User = cfg.user;
          Group = cfg.group;
          Type = "simple";
          ExecStart = "${lib.getExe cfg.package} -config ${configFile}";
          Restart = "always";
        };
      };

      systemd.paths.legit-repos = {
        description = "repositories directory read by legit";
        wantedBy = [ "default.target" ];

        pathConfig = {
          PathExists = cfg.config.repo.scanPath;
          Unit = "legit.service";
        };
      };
    };
}
