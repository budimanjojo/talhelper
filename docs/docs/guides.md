# Guides

## Example talconfig.yaml

A minimal `talconfig.yaml` file will look like this:

```yaml
---
clusterName: my-cluster
endpoint: https://192.168.200.10:6443
nodes:
  - hostname: master
    controlPlane: true
    ipAddress: 192.168.200.11
    installDisk: /dev/sda
```

Let's say you want to add labels to the `master` node and add another worker node named `warmachine`, you can modify `talconfig.yaml` like so:

```{.yaml hl_lines="9-15"}
---
clusterName: my-cluster
endpoint: https://192.168.200.10:6443
nodes:
  - hostname: master
    controlPlane: true
    ipAddress: 192.168.200.11
    installDisk: /dev/sda
    nodeLabels:
      rack: rack1
  - hostname: warmachine
    controlPlane: false
    ipAddress: 192.168.200.12
    installDiskSelector:
      size: 128GB
```

Then you can run `talhelper genconfig`.
Here's a more detailed example [talconfig.yaml](https://github.com/budimanjojo/talhelper/blob/bf3c112f168be583cfc658a5425427974796b2af/example/talconfig.yaml).

To see all the available options of the configuration file, head over to [Configuration Reference](reference/configuration.md).

## Adding Talos extensions and kernel arguments

Talos v1.5 introduced a new unified way to generate boot assets for installer container image that you can build yourself using their `imager` container or use [image-factory](https://factory.talos.dev/) to dynamically build it for you.
The old way of installing system extensions using `machine.install.extensions` in the configuration file is being deprecated, so it's not recommended to use it.

`Talhelper` can help you to generate the installer url like `image-factory` if you provide `schematic` for your nodes.
Let's say your `warmachine` node is using Intel processor so you want to have `intel-ucode` extension and you also want to use traditional network interface naming by providing `net.ifnames=0` to the kernel argument.
Your `talhelper.yaml` should be something like this:

```{.yaml hl_lines="9-15"}
---
clusterName: my-cluster
talosVersion: v1.5.4
endpoint: https://192.168.200.10:6443
nodes:
  - hostname: warmachine
    controlPlane: false
    ipAddress: 192.168.200.12
    schematic:
      customization:
        extraKernelArgs:
          - net.ifnames=0
        systemExtensions:
          officialExtensions:
            - siderolabs/intel-ucode
```

When you run `talhelper genconfig`, the generated manifest file for `warmachine` will have `machine.install.image` value of `factory.talos.dev/installer/9e8cc193609699825d61c039c7738d81cf29c7b20f2a91d8f5c540511b9ea0b4:v1.5.4`, which will be the same url you'll get if using `image-factory`.

If you don't want to use the url from `image-factory` or you want to use your own installer image, you can use per node `talosImageURL` like this:

```{.yaml hl_lines="9-10"}
---
clusterName: my-cluster
talosVersion: v1.5.4
endpoint: https://192.168.200.10:6443
nodes:
  - hostname: warmachine
    controlPlane: false
    ipAddress: 192.168.200.12
    talosImageURL: my.registry/install/talos-installer-image
```

This will result in `machine.install.image` value to be `my.registry/install/talos-installer-image:v1.5.4`.

## Configuring SOPS for Talhelper

[sops](https://github.com/getsops/sops) is a simple and flexible tool for managing secrets.

If you haven't used `sops` before, the easiest way to get started is by using [age](https://github.com/FiloSottile/age) as the encryption tool of choice.
To configure `talhelper` to use `sops` to encrypt and decrypt your secrets, here's the simplified step by step you can do:

1. Install both `sops` and `age` into your system.
2. Run `age-keygen -o <sops-config-dir>/age/keys.txt`. By default, `<sops-config-dir>` will be in `$XDG_CONFIG_HOME/sops` on Linux, `$HOME/Library/Application Support/sops` on MacOS, and `%AppData%\sops` on Windows.
3. In the directory where your `talenv.sops.yaml`, and `talsecrets.sops.yaml` lives, create a `.sops.yaml` file with this content:

    ```yaml
    ---
    creation_rules:
      - age: >-
          <age-public-key> ## get this in the keys.txt file from previous step
    ```

4. Now, if you encrypt your `talenv.sops.yaml` and `talsecret.sops.yaml` files with `sops`, `talhelper` will be able to decrypt it when generating config files.

## Using Doppler instead of SOPS

If you don't want to use `sops` as your secret management, you can use [Doppler](https://www.doppler.com/) instead (or any other secret managers that can inject environment variables to the shell).

Thanks to [@truxnell](https://github.com/truxnell) for this genius idea.
Here's the simplified step by step to achieve this:

1. In the place where you want to use environment secrets, put it in `talconfig.yaml` file with `${}` placeholder like this:

    ```yaml
    controlPlane:
      inlinePatch:
        cluster:
          aescbcEncryptionSecret: ${AESCBCENCYPTIONKEY}
    ```

2. In `doppler`, create a project named i.e "talhelper". In that project, create a config i.e "env" that stores key and value of the secret like `AESCBCENCYPTIONKEY: <secret>.`.
3. Run `doppler` CLI command that sets environment variable before running the `talhelper` command i.e: `doppler run -p talhelper -c env talhelper genconfig`.

Thanks to [@jtcressy](https://github.com/jtcressy) you can also make use of `talsecret.yaml` file (which is a better way than doing `inlinePatch`).
Note that you can only put the cluster secrets known by Talos here (you can use `talhelper gensecret` command and modify it).
Here's the simplified step by step to achieve this:

1. In `talsecret.yaml` file, put all your secrets with `${}` placeholder like this:

    ```yaml
    cluster:
      id: ${CLUSTERNAME}
      secret: ${CLUSTERSECRET}
    secrets:
      bootstraptoken: ${BOOTSTRAPTOKEN}
      secretboxencryptionsecret: ${AESCBCENCYPTIONKEY}
    trustdinfo:
      token: ${TRUSTDTOKEN}
    certs:
      etcd:
        crt: ${ETCDCERT}
        key: ${ETCDKEY}
      k8s:
        crt: ${K8SCERT}
        key: ${K8SKEY}
      k8saggregator:
        crt: ${K8SAGGCERT}
        key: ${K8SAGGKEY}
      k8sserviceaccount:
        key: ${K8SSAKEY}
      os:
        crt: ${OSCERT}
        key: ${OSKEY}
    ```
2. In `doppler`, create a project named i.e "talhelper". In that project, create a config i.e "env" that stores key and value of the secret like `AESCBCENCYPTIONKEY: <secret>.`.
3. Run `doppler` CLI command that sets environment variable before running the `talhelper` command i.e: `doppler run -p talhelper -c env talhelper genconfig`.

## Generating `talosctl` commands for bash scripting

Thanks to the idea and contribution of [mirceanton](https://github.com/mirceanton), you can generate `talosctl` commands for bash scripting in your workflow.
For example, in the directory where a `talconfig.yaml` like this is located:

```yaml
---
clusterName: my-cluster
talosVersion: v1.5.5
nodes:
  - hostname: node1
    ipAddress: 192.168.10.11
    controlPlane: true
```

After running `talhelper genconfig`, you can run `talhelper gencommand apply | bash` in the terminal to apply the generated config into your machine(s) automatically.
There are some other `gencommand` commands that you can use like `upgrade`, `upgrade-k8s`, `bootstrap`, etc,

For more information about the available `gencommand` commands and flags you can use, head over to the [documentation](./reference/cli.md#talhelper-gencommand).

## Selfhosted Image Factory

By default, the generated manifests will use the official [image-factory](https://factory.talos.dev) to pull the installer image.
If you're self hosting your own image-factory, you can change your `talconfig.yaml` like so:

```yaml
---
clusterName: my-cluster
imageFactory:
  registryURL: myfactory.com
  schematicEndpoint: /schematics
  protocol: http
  installerURLTmpl: {{.RegistryURL}}/installer/{{.ID}}:{{.Version}}
```

The `schematicEndpoint` is used to do HTTP POST request to get the schematic ID.
If your selfhosted image factory doesn't do schematic ID like the official one does, you can pass `--offline` flag to `talhelper genconfig` command and modify the `installerURLTmpl` to your needs.

## Editing `talconfig.yaml` file

If you're using a text editor with `yaml` LSP support, you can use `talhelper genschema` command to generate a `talconfig.json`.
You can then feed that file to the language server so you can get autocompletion when editing `talconfig.yaml` file.

If your LSP is configured to use [JSON schema store](https://www.schemastore.org/json/), you should get auto-completion working immediately.

## Shell completion

Depending on how you install `talhelper`, you might not need to do anything to get autocompletion for `talhelper` commands i.e if you install using the Nix Flakes or AUR.

If you don't get it working out of the box, you can use `talhelper completion` command to generate autocompletion for your shell.

=== "bash"
    You will need [bash-completion](https://github.com/scop/bash-completion) installed and configured on your system first.

    And then you can put this line somewhere inside your `~/.bashrc` file:

    ```bash
    source <(talhelper completion bash)
    ```

    After reloading your shell, autocompletion should be working. To enable bash autocompletion in current session of shell, source the `~/.bashrc` file:
    ```bash
    source ~/.bashrc
    ```

=== "fish"
    Put this line somewhere inside your `~/.config/fish/config.fish` file:

    ```fish
    talhelper completion fish | source
    ```

    Another way is to put the generated file into `~/.config/fish/completions/talhelper.fish` file:

    ```fish
    talhelper completion fish > ~/.config/fish/completions/talhelper.fish
    ```

    After reloading your shell, autocompletion should be working.

=== "zsh"
    Put this line somewhere inside your `~/.zshrc`:

    ```zsh
    source <(talhelper completion zsh)
    ```

    After reloading your shell, autocompletion should be working. To enable zsh autocompletion in current session of shell, source the `~/.zshrc` file:
    ```zsh
    source ~/.zshrc
    ```

=== "powershell"
    Append the generated file into `$PROFILE`:

    ```powershell
    talhelper completion powershell >> $PROFILE
    ```

    After reloading your shell, autocompletion should be working.
