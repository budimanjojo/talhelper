# CLI

## talhelper

This command will print a long help introduction to use `talhelper`.

Flags:

```
-h, --help      help for talhelper
-v, --version   show the current talhelper version
```

## talhelper completion

This command will generate autocompletion script for the specified shell.

Usage:

```
talhelper completion [option]
```

Options:

```
bash         Generate autocompletion script for `bash` shell
fish         Generate autocompletion script for `fish` shell
powershell   Generate autocompletion script for `powershell` shell
zsh          Generate autocompletion script for `zsh` shell
```

Flags:

```
-h, --help   help for completion
```

## talhelper genconfig

This command will generate Talos cluster configuration files.

Usage:

```
talhelper genconfig [flag]
```

Flags:

```
-c, --config-file string    File containing configurations for talhelper (default "talconfig.yaml")
-n, --dry-run               Skip generating manifests and show diff instead
-e, --env-file strings      List of files containing env variables for config file (default [talenv.yaml,talenv.sops.yaml,talenv.yml,talenv.sops.yml])
-h, --help                  help for genconfig
    --no-gitignore          Create/update gitignore file too
```

## talhelper gensecret

This command will generate Talos cluster secrets.

Usage:

```
talhelper gensecret [flag]
```

Flags:

```
-f, --from-configfile string   Talos cluster node configuration file to generate secret from
-h, --help                     help for gensecret
```

## talhelper validate

This command will validate the correctness of talconfig or talos node config.

Usage:

```
talhelper validate [option] [file]
```

Option:

```
talconfig     Check the validity of Talhelper config file
nodeconfig    Check the validity of Talos node config file
```

Flags:

```
-h, --help          help for validate
```

Flags available on nodeconfig option:
```
-m, --mode string   Talos runtime mode to validate with (default "metal")
```
