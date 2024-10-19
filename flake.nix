{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "aarch64-linux"
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      imports = [
        inputs.flake-parts.flakeModules.easyOverlay
      ];
      perSystem =
        { config, system, pkgs, ... }:
        {
          overlayAttrs = {
            inherit (config.packages) default;
          };
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            # overlays = [
            #   (final: prev: {
            #     go_1_22 = prev.go_1_22.overrideAttrs (old: {
            #       src = prev.fetchurl {
            #         url = "https://go.dev/dl/go1.22.3.src.tar.gz";
            #         hash = "sha256-gGSO80+QMZPXKlnA3/AZ9fmK4MmqE63gsOy/+ZGnb2g=";
            #       };
            #     });
            #   })
            # ];
          };
          packages = rec {
            default = talhelper;
            talhelper = pkgs.callPackage ./default.nix { };
          };
          devShells.default =
            with pkgs;
            mkShell {
              name = "talhelper-dev";
              packages = [
                gcc
                go_1_23
              ];
            };
        };
    };
}
