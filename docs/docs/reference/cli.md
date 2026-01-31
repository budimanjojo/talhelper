# CLI

## talhelper completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(talhelper completion bash)

To load completions for every new session, execute once:

#### Linux:

	talhelper completion bash > /etc/bash_completion.d/talhelper

#### macOS:

	talhelper completion bash > $(brew --prefix)/etc/bash_completion.d/talhelper

You will need to start a new shell for this setup to take effect.


```
talhelper completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper completion](#talhelper-completion)	 - Generate the autocompletion script for the specified shell

## talhelper completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	talhelper completion fish | source

To load completions for every new session, execute once:

	talhelper completion fish > ~/.config/fish/completions/talhelper.fish

You will need to start a new shell for this setup to take effect.


```
talhelper completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper completion](#talhelper-completion)	 - Generate the autocompletion script for the specified shell

## talhelper completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	talhelper completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
talhelper completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper completion](#talhelper-completion)	 - Generate the autocompletion script for the specified shell

## talhelper completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(talhelper completion zsh)

To load completions for every new session, execute once:

#### Linux:

	talhelper completion zsh > "${fpath[1]}/_talhelper"

#### macOS:

	talhelper completion zsh > $(brew --prefix)/share/zsh/site-functions/_talhelper

You will need to start a new shell for this setup to take effect.


```
talhelper completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper completion](#talhelper-completion)	 - Generate the autocompletion script for the specified shell

## talhelper completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for talhelper for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster
* [talhelper completion bash](#talhelper-completion-bash)	 - Generate the autocompletion script for bash
* [talhelper completion fish](#talhelper-completion-fish)	 - Generate the autocompletion script for fish
* [talhelper completion powershell](#talhelper-completion-powershell)	 - Generate the autocompletion script for powershell
* [talhelper completion zsh](#talhelper-completion-zsh)	 - Generate the autocompletion script for zsh

## talhelper gencommand apply

Generate talosctl apply-config commands.

```
talhelper gencommand apply [flags]
```

### Options

```
  -h, --help   help for apply
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand bootstrap

Generate talosctl bootstrap commands.

```
talhelper gencommand bootstrap [flags]
```

### Options

```
  -h, --help   help for bootstrap
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand health

Generate talosctl health commands.

```
talhelper gencommand health [flags]
```

### Options

```
  -h, --help   help for health
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand kubeconfig

Generate talosctl kubeconfig commands.

```
talhelper gencommand kubeconfig [flags]
```

### Options

```
  -h, --help   help for kubeconfig
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand reset

Generate talosctl reset commands.

```
talhelper gencommand reset [flags]
```

### Options

```
  -h, --help   help for reset
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand upgrade

Generate talosctl upgrade commands.

```
talhelper gencommand upgrade [flags]
```

### Options

```
  -h, --help           help for upgrade
      --offline-mode   Generate schematic ID without doing POST request to image-factory
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand upgrade-k8s

Generate talosctl upgrade-k8s commands.

```
talhelper gencommand upgrade-k8s [flags]
```

### Options

```
  -h, --help   help for upgrade-k8s
```

### Options inherited from parent commands

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -d, --debug                 Whether to enable debugging mode
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### SEE ALSO

* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.

## talhelper gencommand

Generate commands for talosctl.

### Options

```
  -c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
  -e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
      --extra-flags strings   List of additional flags that will be injected into the generated commands.
  -h, --help                  help for gencommand
  -n, --node string           A specific node to generate the command for. If not specified, will generate for all nodes.
  -o, --out-dir string        Directory that contains the generated config files to apply. (default "./clusterconfig")
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster
* [talhelper gencommand apply](#talhelper-gencommand-apply)	 - Generate talosctl apply-config commands.
* [talhelper gencommand bootstrap](#talhelper-gencommand-bootstrap)	 - Generate talosctl bootstrap commands.
* [talhelper gencommand health](#talhelper-gencommand-health)	 - Generate talosctl health commands.
* [talhelper gencommand kubeconfig](#talhelper-gencommand-kubeconfig)	 - Generate talosctl kubeconfig commands.
* [talhelper gencommand reset](#talhelper-gencommand-reset)	 - Generate talosctl reset commands.
* [talhelper gencommand upgrade](#talhelper-gencommand-upgrade)	 - Generate talosctl upgrade commands.
* [talhelper gencommand upgrade-k8s](#talhelper-gencommand-upgrade-k8s)	 - Generate talosctl upgrade-k8s commands.

## talhelper genconfig

Generate Talos cluster config YAML files

```
talhelper genconfig [flags]
```

### Options

```
  -c, --config-file string      File containing configurations for talhelper (default "talconfig.yaml")
      --crt-ttl duration        certificate TTL (default 8760h0m0s)
      --disable-nodes-section   Disable filling the taloscontrol nodes section
  -n, --dry-run                 Skip generating manifests and show diff instead
  -e, --env-file strings        List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -h, --help                    help for genconfig
      --no-gitignore            Create/update gitignore file too
      --offline-mode            Generate schematic ID without doing POST request to image-factory
  -o, --out-dir string          Directory where to dump the generated files (default "./clusterconfig")
  -s, --secret-file strings     List of files containing secrets for the cluster (default [talsecret.yaml,talsecret.sops.yaml,talsecret.yml,talsecret.sops.yml])
  -m, --talos-mode string       Talos runtime mode to validate generated config (default "metal")
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster

## talhelper genschema

Generate `talconfig.yaml` JSON schema file

```
talhelper genschema [flags]
```

### Options

```
  -f, --file string   Where to dump the generated json-schema file (default "talconfig.json")
  -h, --help          help for genschema
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster

## talhelper gensecret

Generate Talos cluster secrets

```
talhelper gensecret [flags]
```

### Options

```
  -f, --from-configfile string   Talos cluster node configuration file to generate secret from
  -h, --help                     help for gensecret
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster

## talhelper genurl image

Generate URL for Talos ISO or disk image

```
talhelper genurl image [flags]
```

### Options

```
  -a, --arch string          CPU architecture support of the image (default "amd64")
      --boot-method string   Boot method of the image (can be disk-image, iso, or pxe) (default "iso")
  -h, --help                 help for image
      --suffix string        The image file extension (only used when boot-method is not iso) (e.g: raw.xz, raw.tar.gz, qcow2)
      --use-uki              Whether to generate UKI image url if Secure Boot is enabled
```

### Options inherited from parent commands

```
  -c, --config-file string          File containing configurations for talhelper (default "talconfig.yaml")
      --customization-file string   File containing customization spec, this will ignore talconfig.yaml file
  -d, --debug                       Whether to enable debugging mode
      --env-file strings            List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -e, --extension strings           Official extension image to be included in the image (ignored when talconfig.yaml is found)
  -k, --kernel-arg strings          Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)
  -n, --node string                 A specific node to generate command for. If not specified, will generate for all nodes (ignored when talconfig.yaml is not found)
      --offline-mode                Generate schematic ID without doing POST request to image-factory
  -r, --registry-url string         Registry url of the image (default "factory.talos.dev")
      --secure-boot                 Whether to generate Secure Boot enabled URL
  -m, --talos-mode string           Talos runtime mode to generate URL (default "metal")
  -v, --version string              Talos version to generate (defaults to latest Talos version) (default "v1.12.2")
```

### SEE ALSO

* [talhelper genurl](#talhelper-genurl)	 - Generate URL for Talos installer or ISO

## talhelper genurl installer

Generate URL for Talos installer image

```
talhelper genurl installer [flags]
```

### Options

```
  -h, --help   help for installer
```

### Options inherited from parent commands

```
  -c, --config-file string          File containing configurations for talhelper (default "talconfig.yaml")
      --customization-file string   File containing customization spec, this will ignore talconfig.yaml file
  -d, --debug                       Whether to enable debugging mode
      --env-file strings            List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -e, --extension strings           Official extension image to be included in the image (ignored when talconfig.yaml is found)
  -k, --kernel-arg strings          Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)
  -n, --node string                 A specific node to generate command for. If not specified, will generate for all nodes (ignored when talconfig.yaml is not found)
      --offline-mode                Generate schematic ID without doing POST request to image-factory
  -r, --registry-url string         Registry url of the image (default "factory.talos.dev")
      --secure-boot                 Whether to generate Secure Boot enabled URL
  -m, --talos-mode string           Talos runtime mode to generate URL (default "metal")
  -v, --version string              Talos version to generate (defaults to latest Talos version) (default "v1.12.2")
```

### SEE ALSO

* [talhelper genurl](#talhelper-genurl)	 - Generate URL for Talos installer or ISO

## talhelper genurl

Generate URL for Talos installer or ISO

### Options

```
  -c, --config-file string          File containing configurations for talhelper (default "talconfig.yaml")
      --customization-file string   File containing customization spec, this will ignore talconfig.yaml file
      --env-file strings            List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -e, --extension strings           Official extension image to be included in the image (ignored when talconfig.yaml is found)
  -h, --help                        help for genurl
  -k, --kernel-arg strings          Kernel arguments to be passed to the image kernel (ignored when talconfig.yaml is found)
  -n, --node string                 A specific node to generate command for. If not specified, will generate for all nodes (ignored when talconfig.yaml is not found)
      --offline-mode                Generate schematic ID without doing POST request to image-factory
  -r, --registry-url string         Registry url of the image (default "factory.talos.dev")
      --secure-boot                 Whether to generate Secure Boot enabled URL
  -m, --talos-mode string           Talos runtime mode to generate URL (default "metal")
  -v, --version string              Talos version to generate (defaults to latest Talos version) (default "v1.12.2")
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster
* [talhelper genurl image](#talhelper-genurl-image)	 - Generate URL for Talos ISO or disk image
* [talhelper genurl installer](#talhelper-genurl-installer)	 - Generate URL for Talos installer image

## talhelper validate nodeconfig

Check the validity of Talos node config file

```
talhelper validate nodeconfig [file] [flags]
```

### Options

```
  -h, --help          help for nodeconfig
  -m, --mode string   Talos runtime mode to validate with (default "metal")
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper validate](#talhelper-validate)	 - Validate the correctness of talconfig or talos node config

## talhelper validate talconfig

Check the validity of talhelper config file

```
talhelper validate talconfig [file] [flags]
```

### Options

```
  -e, --env-file strings   List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
  -h, --help               help for talconfig
      --no-substitute      Whether to do envsubst on before validation
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper validate](#talhelper-validate)	 - Validate the correctness of talconfig or talos node config

## talhelper validate

Validate the correctness of talconfig or talos node config

### Options

```
  -h, --help   help for validate
```

### Options inherited from parent commands

```
  -d, --debug   Whether to enable debugging mode
```

### SEE ALSO

* [talhelper](#talhelper)	 - A tool to help with creating Talos cluster
* [talhelper validate nodeconfig](#talhelper-validate-nodeconfig)	 - Check the validity of Talos node config file
* [talhelper validate talconfig](#talhelper-validate-talconfig)	 - Check the validity of talhelper config file

## talhelper

A tool to help with creating Talos cluster

### Synopsis

talhelper is a tool to help you create a Talos cluster.

Workflow:
  Create talconfig.yaml file defining your nodes information like so:

```
  clusterName: mycluster
  talosVersion: v1.0
  endpoint: https://192.168.200.10:6443
  nodes:
    - hostname: master1
      ipAddress: 192.168.200.11
      installDisk: /dev/sdb
      controlPlane: true
    - hostname: worker1
      ipAddress: 192.168.200.21
      installDisk: /dev/nvme1
      controlPlane: false

```

  Then run these commands:
  > talhelper gensecret > talsecret.sops.yaml
  > sops -e -i talsecret.sops.yaml
  > talhelper genconfig

  The generated yaml files will be in ./clusterconfig directory

  WARNING! Please don't push the generated files into your public git repository.
  By default talhelper will create a ".gitignore" file to ignore the generated files for you unless you use "--no-gitignore" flag.
  The generated files contain unencrypted secrets and you don't want people to get a hand of them.

### Options

```
  -d, --debug   Whether to enable debugging mode
  -h, --help    help for talhelper
```

### SEE ALSO

* [talhelper completion](#talhelper-completion)	 - Generate the autocompletion script for the specified shell
* [talhelper gencommand](#talhelper-gencommand)	 - Generate commands for talosctl.
* [talhelper genconfig](#talhelper-genconfig)	 - Generate Talos cluster config YAML files
* [talhelper genschema](#talhelper-genschema)	 - Generate `talconfig.yaml` JSON schema file
* [talhelper gensecret](#talhelper-gensecret)	 - Generate Talos cluster secrets
* [talhelper genurl](#talhelper-genurl)	 - Generate URL for Talos installer or ISO
* [talhelper validate](#talhelper-validate)	 - Validate the correctness of talconfig or talos node config