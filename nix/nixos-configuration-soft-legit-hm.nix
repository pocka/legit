# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  config,
  lib,
  pkgs,
  ...
}:
let
  home-manager = builtins.fetchTarball {
    url = "https://github.com/nix-community/home-manager/archive/4ce190229c73d44536caa7072f6308fb2d8feeb3.tar.gz";
    sha256 = "1cqangi17i4nfkjpzpzpsavcgnaqdvarjjsjg09sbj5mnhnv6v35";
  };
in
{
  imports = [ (import "${home-manager}/nixos") ];

  system.stateVersion = "26.11";

  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  networking.firewall.enable = false;

  users.users = {
    alice = {
      isNormalUser = true;
      password = "pass";
    };
  };

  security.polkit = {
    enable = true;

    # Password-less poweroff / reboot
    extraConfig = ''
      polkit.addRule(function(action, subject) {
        const isPowerOff = [
          "org.freedesktop.login1.reboot",
          "org.freedesktop.login1.reboot-multiple-sessions",
          "org.freedesktop.login1.power-off",
          "org.freedesktop.login1.power-off-multiple-sessions",
        ].indexOf(action.id) > -1;

        if (subject.user === "alice" && isPowerOff) {
          return polkit.Result.YES;
        }
      })
    '';
  };

  services.getty.autologinUser = "alice";

  services.caddy = {
    enable = true;
    configFile = pkgs.writeText "Caddyfile" ''
      {
        auto_https off
      }

      :80 {
        reverse_proxy :5555
      }
    '';
  };

  home-manager.useUserPackages = true;

  home-manager.users.alice =
    {
      config,
      lib,
      pkgs,
      legit,
      ...
    }:
    let
      write-initial-admin-keys = pkgs.writeShellApplication {
        name = "write-initial-admin-keys";

        text = ''
          echo "SOFT_SERVE_INITIAL_ADMIN_KEYS=\"$(cat)\"" > ${config.xdg.configHome}/soft-serve/admin.env
        '';
      };
    in
    {
      imports = [ legit.homeManagerModules.default ];

      home.stateVersion = "26.05";
      # This is inevitable as we use fetched tarball along with nixpkgs managed by Flake.
      home.enableNixpkgsReleaseCheck = false;
      home.packages = with pkgs; [
        curl
        write-initial-admin-keys
      ];

      services.legit = {
        enable = true;
        config = {
          repo.scanPath = "${config.xdg.dataHome}/soft-serve/repos";
        };
      };

      xdg.configFile."soft-serve/.keep" = {
        text = "";
      };

      xdg.dataFile."soft-serve/config.yaml" =
        let
          yaml = pkgs.formats.yaml { };
        in
        {
          source = yaml.generate "soft-serve.yaml" {
            name = "soft-serve & legit HM module demo";
            log_format = "text";
            ssh.listen_addr = ":22222";
            http.listen_addr = ":8080";
          };
        };

      systemd.user.paths.soft-serve-adminenv = {
        Unit = {
          Description = "soft-serve initial admin env file";
        };

        Install = {
          WantedBy = [ "default.target" ];
        };

        Path = {
          PathExists = "${config.xdg.configHome}/soft-serve/admin.env";
          Unit = "soft-serve.service";
        };
      };

      systemd.user.services.soft-serve = {
        Unit = {
          Description = "Git server for the command-line";
          After = [ "network.target" ];
        };

        Service = {
          Type = "simple";
          Restart = "always";
          RestartSec = 1;
          ExecStart = "${lib.getExe pkgs.soft-serve} serve";
          Environment = "SOFT_SERVE_DATA_PATH=${config.xdg.dataHome}/soft-serve";
          EnvironmentFile = "${config.xdg.configHome}/soft-serve/admin.env";
          WorkingDirectory = "${config.xdg.dataHome}/soft-serve";
        };
      };
    };
}
