# Introduction

## Overview

`talhelper` is a tool to help creating [Talos](https://talos.dev) configuration files declaratively.
It was inspired by a `python` script written by [@bjw-s](https://github.com/bjw-s).
You can say `talhelper` is like `kustomize` but for Talos manifest files with `SOPS` support natively.

In a nutshell, this is what `talhelper` does step by step behind the door:

* Read and validate `talconfig.yaml`.
* Read and decrypt `talsecret.yaml` or `talsecret.sops.yaml` with `sops` if needed.
* Read and decrypt `talenv.yaml` or `talenv.sops.yaml` with `sops` if needed and load them into environment variables.
* Do [envsubst](https://linux.die.net/man/1/envsubst) if needed.
* Validate and generate Talos and machine config files inside `./clusterconfig` directory.
* Generate `.gitignore` file so you don't commit the generated files to the public.

## Why should I use Talhelper

The main reason to use `talhelper` instead of `talosctl gen config` command to generate Talos `machineconfig` files is because you want to have them version controlled in your git repository which is currently not possible yet with `talosctl`.

Currently, to create `Talos` configuration files using the official `talosctl` tool your steps are:

- Run `talosctl gen config <cluster-name> <cluster-endpoint>` and it will generate `controlplane.yaml`, `worker.yaml`, `talosconfig` in the current working directory.
- Copy and modify those files manually according to your nodes.
- Run `talosctl apply-config --insecure -n <ip-address> --file <your-modified-file.yaml>` for each node.

This process is fine if you just want to do this once and forget about it. But if you're like me (and many [others](https://discord.com/invite/home-operations)), you might want to "GitOpsified" this process. So here's where you might want to use `talhelper`.

With `talhelper`, the steps will become like this:

- Create a `talconfig.yaml`.
- Run `talhelper gensecret > talsecret.sops.yaml` and encrypt it with [sops](https://github.com/getsops/sops) `sops -e -i talsecret.sops.yaml`.
- Run `talhelper genconfig`.
- Run `talosctl apply-config --insecure -n <ip-address> --file ./clusterconfig/<cluster-name>-<hostname>.yaml` for each node.

Yes there are more steps needed.
But now you can commit your `talconfig.yaml` and the encrypted `talsecret.sops.yaml` to your repository and get your whole cluster version controlled.

To get started, hop over to the [Getting Started](getting-started.md) section.

## Alternatives

There are some alternatives you can consider instead of `talhelper`.

- The official [Terraform provider](https://registry.terraform.io/providers/siderolabs/talos/latest)
- The official [Pulumi provider](https://www.pulumi.com/registry/packages/talos/)
- [TOPF](https://github.com/postfinance/topf) — manages Talos cluster lifecycle with layered patches and SOPS support

## Bug report and feature request

If you have encountered any bug or you want to request a new feature, please [open an issue](https://github.com/budimanjojo/talhelper/issues/new) at GitHub.
