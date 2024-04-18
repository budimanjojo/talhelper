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
      perSystem =
        { system, pkgs, ... }:
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [
              (final: prev: {
                go_1_22 = prev.go_1_22.overrideAttrs (old: {
                  src = prev.fetchurl {
                    url = "https://go.dev/dl/go1.22.2.src.tar.gz";
                    hash = "sha256-N06oKyiexzjpaCZ8rFnH1f8YD5SSJQJUeEsgROkN9ak=";
                  };
                });
              })
            ];
          };
          packages.default = pkgs.callPackage ./default.nix { };
          devShells.default =
            with pkgs;
            mkShell {
              name = "talhelper-dev";
              packages = [
                gcc
                go_1_22
              ];
            };
        };
    };
}
