<div align="center">
  <h3 align="center">Talhelper</h3>

  ![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/budimanjojo/talhelper?include_prereleases)
  ![GitHub issues](https://img.shields.io/github/issues/budimanjojo/talhelper)
  ![GitHub](https://img.shields.io/github/license/budimanjojo/talhelper)

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
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

## About The Project

The main reason of this tool is to help creating Talos cluster in GitOps way.
Inspired by a python script written by [@bjw-s](https://github.com/bjw-s) [here](https://github.com/bjw-s/home-ops/blob/main/infrastructure/talos/buildClusterConfig.py).

This tool will:
* Read your `talconfig.yaml`
* Read and decrypt your `talenv.yaml` with [SOPS](https://github.com/mozilla/sops)
* Do [envsubst](https://linux.die.net/man/1/envsubst) if needed
* Validate config file is good for talosctl
* Generate Talos cluster and config yaml files for you based on your `talconfig.yaml`
* Generate `.gitignore` file so you don't commit your secret to the public

This tool is actually my first time programming something other than shell script.
Any input and suggestion will be highly appreciated.

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting Started

1. Create a `talconfig.yaml`, an example [template](./test/talconfig.yaml) is provided.
2. Run `talhelper gensecret --patch-configfile > talenv.yaml` (`--patch-configfile` will add inlinePatches inside your `talconfig.yaml`)
3. Encrypt the secret with SOPS: `sops -e -i talenv.yaml`
4. Run `talhelper genconfig` and the output files will be in `./clusterconfig` by default.

To get help, run `talhelper <subcommand> --help`

### Installation

TBD

<p align="right">(<a href="#top">back to top</a>)</p>

## Usage

```
Available Commands:
  completion  Generate the autocompletion script for the specified shell
  genconfig   Generate Talos cluster config YAML file
  gensecret   Generate Talos cluster secrets
  help        Help about any command
```

```
  talhelper genconfig [flags]

Flags:
  -c, --config-file string   File containing configurations for nodes (default "talconfig.yaml")
  -e, --env-file string      File containing env variables for config file (default "talenv.yaml")
  -h, --help                 help for genconfig
      --no-gitignore         Create/update gitignore file too
  -o, --out-dir string       Directory where to dump the generated files (default "./clusterconfig")
```

```
Usage:
  talhelper gensecret [flags]

Flags:
  -c, --config-file string   File containing configurations for talhelper (default "talconfig.yaml")
  -h, --help                 help for gensecret
  -p, --patch-configfile     Whether to generate inline patches into config file
```
<p align="right">(<a href="#top">back to top</a>)</p>

## Roadmap

- [ ] Add tests
- [ ] Add release workflows
- [ ] More useful features

See the [open issues](https://github.com/othneildrew/Best-README-Template/issues) for a full list of proposed features (and known issues).

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
