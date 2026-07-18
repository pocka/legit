# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  config,
  lib,
  pkgs,
  ...
}:
let
  # Directory managed by soft-serve.
  softserveData = "/opt/soft";
  softserveEnvDir = "/opt/soft-admin";
  softserveAdminKeyEnv = "${softserveEnvDir}/admin.env";

  write-initial-admin-keys = pkgs.writeShellApplication {
    name = "write-initial-admin-keys";

    text = ''
      echo "SOFT_SERVE_INITIAL_ADMIN_KEYS=\"$(cat)\"" > ${softserveAdminKeyEnv}
    '';
  };
in
{
  system.stateVersion = "26.11";

  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  networking.firewall.enable = false;

  users.groups.git = { };

  systemd.tmpfiles.rules = [
    "d ${softserveData} 0750 soft git -"
    "d ${softserveEnvDir} 0770 soft git -"
  ];

  users.users = {
    soft = {
      isSystemUser = true;
      group = "git";
    };
    legit = {
      isSystemUser = true;
      group = "git";
    };
    admin = {
      isNormalUser = true;
      password = "pass";
      group = "git";
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

        if (subject.user === "admin" && isPowerOff) {
          return polkit.Result.YES;
        }
      })
    '';
  };

  services.getty.autologinUser = "admin";

  environment.systemPackages = with pkgs; [
    curl
    soft-serve
    write-initial-admin-keys
  ];

  systemd.paths.soft-serve-adminenv = {
    description = "availability of initial admin env file";
    wantedBy = [ "default.target" ];

    pathConfig = {
      PathExists = softserveAdminKeyEnv;
      Unit = "soft-serve.service";
    };
  };

  # Service definition in nixpkgs isn't good.
  systemd.services.soft-serve =
    let
      yaml = pkgs.formats.yaml { };
      configFile = yaml.generate "soft-serve.yaml" {
        name = "soft-serve & legit demo";
        log_format = "text";
        ssh.listen_addr = ":22222";
        http.listen_addr = ":8080";
      };
    in
    {
      description = "Git server for the command-line";
      after = [ "network.target" ];
      restartTriggers = [
        configFile
        softserveAdminKeyEnv
      ];
      wantedBy = lib.mkForce [ ];

      environment = {
        SOFT_SERVE_DATA_PATH = softserveData;
        SOFT_SERVE_CONFIG_LOCATION = configFile;
      };

      serviceConfig = {
        User = "soft";
        Group = "git";
        Type = "simple";
        ExecStart = "${lib.getExe pkgs.soft-serve} serve";
        EnvironmentFile = [ softserveAdminKeyEnv ];
        WorkingDirectory = softserveData;
      };
    };

  services.legit = {
    enable = true;
    config = {
      repo.scanPath = "${softserveData}/repos";
    };
    user = "legit";
    group = "git";
  };

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
}
