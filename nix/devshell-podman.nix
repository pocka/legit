# Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
# SPDX-License-Identifier: MIT

{
  stdenv,
  mkShell,
  podman,
}:

let
  conf = stdenv.mkDerivation {
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
mkShell {
  packages = [
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
}
