#!/usr/bin/env bash

ROOTDIR=$(git rev-parse --show-toplevel)
NIXFILE="${ROOTDIR}"/default.nix
LAST_COMMIT=$(git rev-parse HEAD)

setKV () {
  if [ "$2" != "" ]; then
    echo "setting $1 to $2"
  fi
  sed -i "s|$1 = \".*\"|$1 = \"${2:-}\"|" "${NIXFILE}"
}

hash=$(nix-prefetch-url --quiet --unpack https://github.com/budimanjojo/talhelper/archive/"${LAST_COMMIT}".tar.gz)
SHA256=$(nix hash to-sri --type sha256 "${hash}")

setKV version "${LAST_COMMIT}"
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

echo "try building with commit ${LAST_COMMIT}..."
nix build
