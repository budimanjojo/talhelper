# Getting Started

## Before you begin

There are some prerequisites before you start using `talhelper`.

- You need `talhelper` installed on your workstation (of course), head over to the [Installation](installation.md) page for more detail.
- You also need `sops` installed and configured with your preferred encryption tool (`age`, `pgp`, etc). If you want to use `doppler` instead, there's an alternative way to do this thanks to [@truxnell](https://github.com/truxnell) which you can read [here](guides.md#using-doppler-instead-of-sops)
- You also need `talosctl` installed on your workstation to apply the generated machine config files.

Once you have all of the above conditions met, you can now start with the [Scenario](#scenarios) that suits your current situation.

## Scenarios

Depending on which situation you are currently in before integrating `talhelper` to your stack, here are some simplified steps to get you started:

### You already have a Talos cluster running

If you already have your Talos Kubernetes cluster up and running but you haven't GitOps it yet.
Here are the steps you need to do:

1. Get your node's `machineconfig` using `talosctl`: `talosctl -n <node-ip> read /system/state/config.yaml > /tmp/machineconfig.yaml`.
2. Run `talhelper gensecret -f /tmp/machineconfig.yaml > talsecret.sops.yaml`. This command will create a `talsecret.sops.yaml` file with all your current cluster secrets.
3. Encrypt the secret with `sops`: `sops -e -i talsecret.sops.yaml` (you will need `sops` [configured properly](guides.md#configuring-sops-for-talhelper)).
4. Create a `talconfig.yaml` based on your current cluster, here's the example [template](https://github.com/budimanjojo/talhelper/blob/master/example/talconfig.yaml). For all the available options, look at the [Configuration Reference](reference/configuration.md)
5. Run `talhelper genconfig` and the output files will be in `./clusterconfig` by default.
6. Commit your `talconfig.yaml` and `talsecret.yaml` in your git repository.

!!! note

    Please don't push the generated files into your public git repository.
    By default `talhelper` will create a `.gitignore` file to ignore the generated files for you unless you use `--no-gitignore` flag.

    The generated files contain unencrypted secrets and you don't want people to get a hand on them.

### You are starting from scratch

If you are creating a Talos Kubernetes cluster from scratch and you want to use `talhelper`, that's awesome!
Here are the steps you need to do:

1. Create a `talconfig.yaml` according to your needs, here's the example [template](https://github.com/budimanjojo/talhelper/blob/master/example/talconfig.yaml). For all the available options, look at the [Configuration Reference](reference/configuration.md)
2. Run `talhelper gensecret > talsecret.sops.yaml`. This command will create a `talsecret.sops.yaml` file with your future cluster secrets.
3. Encrypt the secret with `sops`: `sops -e -i talsecret.sops.yaml` (you will need `sops` [configured properly](guides.md#configuring-sops-for-talhelper)).
4. Run `talhelper genconfig` and the output files will be in `./clusterconfig` by default.
5. Commit your `talconfig.yaml` and `talsecret.yaml` in your git repository.

!!! note

    Please don't push the generated files into your public git repository.
    By default `talhelper` will create a `.gitignore` file to ignore the generated files for you unless you use `--no-gitignore` flag.

    The generated files contain unencrypted secrets and you don't want people to get a hand on them.
