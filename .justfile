set shell := ['bash', '-euo', 'pipefail', '-c']

# `nix` is required
requireNixToRun := require('nix')

export NIX_CONFIG := 'extra-experimental-features = nix-command flakes'
export SOPS_AGE_KEY := 'AGE-SECRET-KEY-172FENV3SDP8JSRRX2SWTA9JQMAW7MW3GSKJ2JZDNXS4GVFAS5STQUW8WN4'

# we use nix shebang to get into a nix shell with packages that's needed for each recipe
nixShebang := '/usr/bin/env -S nix shell --inputs-from ' + justfile_directory()
realShebang := '/usr/bin/env bash -euo pipefail'
# we can also get into devShell and run stuffs
nixDevelopShebang := '/usr/bin/env -S nix develop -c ' + realShebang

[private]
default:
    @just -l

# deploy docs site to be accessed locally
[working-directory: 'docs']
deploy-docs-site port='1313':
    #! {{ nixDevelopShebang }}
    mkdocs serve -a 0.0.0.0:{{ port }}

# generate CLI reference documentation
gen-cli-ref-docs outdir='docs/docs/reference':
    #! {{ nixDevelopShebang }}
    go run main.go gendocs {{ outdir }}

# generate `talconfig.yaml` JSON schema
gen-config-schema out='pkg/config/schemas/talconfig.json':
    #! {{ nixDevelopShebang }}
    go run main.go genschema -f {{ out }}

# generate `talos-extensions.json` file
[working-directory: 'hack/tsehelper']
gen-extensions-schema out='../../pkg/config/schemas/talos-extensions.json':
    #! {{ nixDevelopShebang }}
    go run . -minimal --output ../../pkg/config/schemas/talos-extensions.json

# run all tests
test: test-build test-coverage lint

# run `go test`
test-coverage:
    #! {{ nixDevelopShebang }}
    go test -v ./... -race -covermode=atomic

# test if the package build and try running it against example
test-build:
    #! {{ nixDevelopShebang }}
    tempdir=$(mktemp -d)
    go build -o "$tempdir"/talhelper
    cd example
    "$tempdir"/talhelper genconfig
    rm -rf "$tempdir"

# run `golangci-lint`
lint:
    #! {{ nixShebang }} nixpkgs#golangci-lint -c {{ realShebang }}
    golangci-lint run --timeout 3m0s

# release new version on GitHub
release-new-version *version:
    #! {{ nixShebang }} nixpkgs#gh -c {{ realShebang }}
    gh release create {{ version }}
