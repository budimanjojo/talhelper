#!/usr/bin/env bash

ROOTDIR=$(git rev-parse --show-toplevel)
NIXFILE="${ROOTDIR}"/default.nix

nixExpr="with import <nixpkgs> {}; let flake = builtins.getFlake "path:${ROOTDIR}"; in flake.outputs.packages.x86_64-linux.default.version"

LATEST=$(curl -s https://api.github.com/repos/budimanjojo/talhelper/releases/latest | jq -r .name)
CURRENT=$(nix-instantiate --eval -E "${nixExpr}" | tr -d '"')

if [ "${LATEST#v}" != "${CURRENT#v}" ]; then
  echo "updating version in $(basename "${NIXFILE}") from ${CURRENT#v} to ${LATEST#v}..."
  setKV () {
    if [ "$2" != "" ]; then
      echo "setting $1 to $2"
    fi
    sed -i "s|$1 = \".*\"|$1 = \"${2:-}\"|" "${NIXFILE}"
  }
  hash=$(nix-prefetch-url --quiet --unpack https://github.com/budimanjojo/talhelper/archive/refs/tags/"${LATEST}".tar.gz)
  SHA256=$(nix hash to-sri --type sha256 "${hash}")
  setKV version "${LATEST#v}"
  setKV sha256 "${SHA256}"
  setKV vendorHash "" # so that the build will fail and provide the actual hash

  set +e
  vendorHash=$(nix build --no-link 2>&1 >/dev/null | grep "got:" | cut -d':' -f2 | sed 's| ||g')
  VENDOR_SHA256=$(nix hash to-sri --type sha256 "${vendorHash}")
  set -e

  if [ -n "${VENDOR_SHA256:-}" ]; then
    setKV vendorHash "${VENDOR_SHA256}"
  else
    echo "Update failed. VENDOR_SHA256 is empty"
    exit 1
  fi

  echo "done updating $(basename "${NIXFILE}"), you can commit the changes now"
else
  echo "version in $(basename "${NIXFILE}") is already the latest ${CURRENT}. nothing to do"
fi
