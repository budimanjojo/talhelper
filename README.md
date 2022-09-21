<div align="center">
  <h3 align="center">Talhelper</h3>

  [![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/budimanjojo/talhelper?include_prereleases)](https://github.com/budimanjojo/talhelper/releases)
  [![GitHub issues](https://img.shields.io/github/issues/budimanjojo/talhelper)](https://github.com/budimanjojo/talhelper/issues)
  [![License](https://img.shields.io/github/license/budimanjojo/talhelper)](https://github.com/budimanjojo/talhelper/blob/master/LICENSE)
  [![AUR link](https://img.shields.io/aur/version/talhelper-bin)](https://aur.archlinux.org/packages/talhelper-bin)

  <p align="center">
    A helper tool to help creating Talos cluster in your GitOps repository.
    <br />
    ·
    <a href="https://github.com/budimanjojo/talhelper/issues">Report Bug</a>
    ·
    <a href="https://github.com/budimanjojo/talhelper/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About The Project

The main reason of this tool is to help creating Talos cluster in GitOps way.
Inspired by a python script written by [@bjw-s](https://github.com/bjw-s).

You can use this tool to generate Talos config file with `talhelper genconfig` command.
You can also use this tool to generate Talos secrets with `talhelper gensecret` command.

This tool will:
* Read your `talconfig.yaml`
* Read and decrypt your `talsecret.yaml` or `talsecret.sops.yaml` with [SOPS](https://github.com/mozilla/sops)
* Read and decrypt your `talenv.yaml` or `talenv.sops.yaml` with [SOPS](https://github.com/mozilla/sops)
* Do [envsubst](https://linux.die.net/man/1/envsubst) if needed
* Validate config file is good for talosctl
* Generate Talos cluster and config yaml files for you based on your `talconfig.yaml`
* Generate `.gitignore` file so you don't commit your secret to the public

This tool is my first time programming something other than shell script.
Any input and suggestion will be highly appreciated.

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

Scenario 1 (You already have your talos config but not GitOps it yet):
1. Create a `talconfig.yaml` based on your current cluster, an example [template](./test/talconfig.yaml) is provided.
2. Run `talhelper gensecret -f <your-talos-controlplane.yaml> > talsecret.sops.yaml`. This will create a `talsecret.sops.yaml` file with all your current cluster secrets.
3. Encrypt the secret with SOPS: `sops -e -i talsecret.sops.yaml`.
4. Run `talhelper genconfig` and the output files will be in `./clusterconfig` by default. Make sure the generated files are identical with your current machine config files.
5. Commit your `talconfig.yaml` and `talsecret.sops.yaml` in Git repository.

Scenario 2 (You want talhelper to create from scratch):
1. Create a `talconfig.yaml`, an example [template](./test/talconfig.yaml) is provided.
2. Run `talhelper gensecret > talsecret.sops.yaml`.
3. Encrypt the secret with SOPS: `sops -e -i talsecret.sops.yaml`.
4. Run `talhelper genconfig` and the output files will be in `./clusterconfig` by default.
5. Commit your `talconfig.yaml` and `talenv.sops.yaml` in Git repository.

To get help, run `talhelper <subcommand> --help`

### Installation

There are several ways to install `talhelper`:
- Using [aqua](https://aquaproj.github.io/).
- Download the archives from [release](https://github.com/budimanjojo/talhelper/releases/latest) page.
- From [AUR](https://aur.archlinux.org/packages/talhelper-bin) for Arch Linux users.
- Install it using this one liner, using tool from [jpillora](https://github.com/jpillora/installer):
  ```
  curl https://i.jpillora.com/budimanjojo/talhelper! | sudo bash
  ```
<p align="right">(<a href="#top">back to top</a>)</p>

## Usage

```
Available Commands:
  completion  Generate the autocompletion script for the specified shell
  genconfig   Generate Talos cluster config YAML files
  gensecret   Generate Talos cluster secrets
  help        Help about any command
  validate    Validate the correctness of talconfig or talos node config
```

```
  talhelper genconfig [flags]

Flags:
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -h, --help                  help for genconfig
      --no-gitignore          Create/update gitignore file too
  -o, --out-dir string        Directory where to dump the generated files (default "./clusterconfig")
  -s, --secret-file strings   List of files containing secrets for the cluster (default [talsecret.yaml,talsecret.sops.yaml,talsecret.yml,talsecret.sops.yml])
  -m, --talos-mode string     Talos runtime mode to validate generated config (default "metal")
```

```
Usage:
  talhelper gensecret [flags]

Flags:
  -f, --from-configfile string   Talos cluster node configuration file to generate secret from
  -h, --help                     help for gensecret
```

```
Usage:
  talhelper validate nodeconfig [file] [flags]

Flags:
  -h, --help          help for nodeconfig
  -m, --mode string   Talos runtime mode to validate with (default "metal")
```

```
Usage:
  talhelper validate talconfig [file] [flags]

Flags:
  -h, --help   help for talconfig
```

<p align="right">(<a href="#top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>

## License

Distributed under the BSD-3 License. See [LICENSE](./LICENSE) for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

## Acknowledgments

* [bjw-s](https://github.com/bjw-s) <- The guy who inspired this tool
* [k8s@home](https://github.com/k8s-at-home/) <- Best community of people running Kubernetes at home
* [Best-README-Template](https://github.com/othneildrew/Best-README-Template) <- Where this README is built off from

<p align="right">(<a href="#top">back to top</a>)</p>
