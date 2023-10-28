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

## Editing `talconfig.yaml` file

If you're using a text editor with `yaml` LSP support, you can use `talhelper genschema` command to generate a `talconfig.json`.
You can then feed that file to the language server so you can get autocompletion when editing `talconfig.yaml` file.

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
