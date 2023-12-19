# tsehelper - talos system extensions helper

This is a tool to generate a system extension to talos version mapping schema from a given talos system extension container registry.

- [tsehelper - talos system extensions helper](#tsehelper---talos-system-extensions-helper)
  - [Usage](#usage)
    - [Environment Variables](#environment-variables)
  - [Usage Examples](#usage-examples)

## Usage

The tool is idempotent, insofar that it will only pull versions that are not already cached. If you'd like to force a refresh already cached versions, you can delete the cache file (or the versions you want to remove from the file) and run the tool again.

The tool will output the generated schema to stdout and to the configured cache file location.

### Environment Variables

There are four environment variables that can be used to configure the tool:

- TSEHELPER_REGEX_OVERRIDE: The container registry to pull the talos system extensions from. Defaults to `ghcr.io/siderolabs/extensions`.
- TSEHELPER_TALOS_EXTENSIONS_REPO: The directory to cache the talos system extensions in. Defaults to the various cache directories provided by the [os.UserCacheDir](https://pkg.go.dev/os#UserCacheDir) function.
- TSEHELPER_TALOS_EXTENSIONS_CACHE_DIR: The filename to write the generated system extensions cache file to. Defaults to `talos-extensions.json`.
- TSEHELPER_TALOS_EXTENSIONS_CACHE_FILE: Override the regex used to parse the talos system extension container tags. Defaults to `^(?P<registry>[\w\.\-0-9]+)\/(?P<org>[\w\.\-0-9]+)\/(?P<repo>[\w\.\-0-9]+):(?P<tag>[\w\.\-0-9]+)@sha256:(?P<shasum>[a-f0-9]+)$`.

There are seven flags that can be used to change the output format:

- `-minimal`: Flag to indicate whether to output minimal json consisting of only the org/repo, e.g. `siderolabs/amd-ucode`. Minimal is mutually exclusive from `-trimRegistry` and `-trimSha256`.

- `-purgecache`: Flag to indicate whether to purge the cache file before generating the schema.

- `-onlyVersions`: Flag to indicate whether to only output the versions, e.g. `v1.5.5`.
- `-version <version>`: Return the system extensions for a specific version, e.g. `v1.4.1`. If not specified, all versions will be returned.

- `-trimRegistry`: Flag to indicate whether to trim the registry prefix, e.g. `ghcr.io/siderolabs/extensions/siderolabs/amd-ucode:v1.2.0@sha256:1234567890abcdef` -> `siderolabs/amd-ucode:v1.2.0@sha256:1234567890abcdef`
- `-trimSha256`: Flag to indicate whether to trim the sha256 suffix e.g. `ghcr.io/siderolabs/amd-ucode:v1.2.0@sha256:1234567890abcdef` -> `ghcr.io/siderolabs/amd-ucode:v1.2.0`
- `-trimTag`: Flag to indicate whether to trim the tag e.g. `ghcr.io/siderolabs/amd-ucode:v1.2.0@sha256:1234567890abcdef` -> `ghcr.io/siderolabs/amd-ucode@sha256:1234567890abcdef

## Usage Examples

To output images by shasum, for instance:

`$ tsehelper -trimTag`

```yaml
{
    "versions": [
        {
            "version": "v1.2.0",
            "systemExtensions": [
                "ghcr.io/siderolabs/amd-ucode@sha256:fbfc7011ea810c485e1351ec585f43dfd2f509d246
ec6d05bea87f8951e5b220",
                "ghcr.io/siderolabs/bnx2-bnx2x@sha256:dc24666533e69aa8a245c8f1920340ef78edae50c
0805dab8660bce0a06381bd",
                "ghcr.io/siderolabs/gvisor@sha256:4391ae38e33d7b27a7d9a28fc120ac03a7f8de53a953c
fe412501bc73005f2fb",
                "ghcr.io/siderolabs/hello-world-service@sha256:154575af40453064e7a89d3ddee3c084
e648fbb1fc40f74d2d7517403cfe9c5e",
                "ghcr.io/siderolabs/intel-ucode@sha256:4f6eda2a0484cd8ec68df37c0e3589c4e5c1970b
0a9e82d6859f92943e7f4624",
                "ghcr.io/siderolabs/nvidia-container-toolkit@sha256:2c909f41c4c6115ede86856a5af
3844035cb75a571649a7510d9b235492e6f49",
                "ghcr.io/siderolabs/nvidia-fabricmanager@sha256:56c933b1f3134ff0dc0426a756339b8
1970b9874e7bf53b3f899d8e810cd529c",
                "ghcr.io/siderolabs/nvidia-open-gpu-kernel-modules@sha256:c3f97c915759d6b569ab8
679603876174d75845495ebd921c7b40179334873eb"
            ]
        }
    ]
}
```

For instance to write a file to the cache the trims the registry prefix and the sha256 suffix, you can run:

`$ tsehelper -trimRegistry -trimSha256`

```yaml
{
    "versions": [
        {
            "version": "v1.2.0",
            "systemExtensions": [
                "siderolabs/amd-ucode:20220411",
                "siderolabs/bnx2-bnx2x:20220411",
                "siderolabs/gvisor:20220405.0-v1.2.0",
                "siderolabs/hello-world-service:v1.2.0",
                "siderolabs/intel-ucode:20220809",
                "siderolabs/nvidia-container-toolkit:515.65.01-v1.10.0",
                "siderolabs/nvidia-fabricmanager:515.65.01",
                "siderolabs/nvidia-open-gpu-kernel-modules:515.65.01-v1.2.0"
            ]
        }
    ]
}
```

`$ tsehelper -minimal`

```yaml
{
    "versions": [
        {
            "version": "v1.2.0",
            "systemExtensions": [
                "siderolabs/amd-ucode,
                "siderolabs/bnx2-bnx2x,
                "siderolabs/gvisor,
                "siderolabs/hello-world-service,
                "siderolabs/intel-ucode,
                "siderolabs/nvidia-container-toolkit,
                "siderolabs/nvidia-fabricmanager,
                "siderolabs/nvidia-open-gpu-kernel-modules"
            ]
        }
    ]
}
```
