# Installation

There are several ways to install `talhelper` to your workstation.

## Using aqua

You can get `talhelper` from the standard registry as `budimanjojo/talhelper`.

## Using asdf

You can get `talhelper` with a [plugin](https://github.com/bjw-s/asdf-talhelper) maintained by [@bjw-s](https://github.com/bjw-s).

- Add the plugin

    ```bash
    asdf plugin add talhelper
    ```

- Install the program

    ```bash
    asdf install talhelper latest
    ```

## Using Homebrew

You can get `talhelper` from the official formulae (thanks to [@ishioni](https://github.com/ishioni)).

```
brew install talhelper
```

## Using Nix Flakes

You can get `talhelper` as [Nix Flakes](https://nixos.wiki/wiki/Flakes) from the [repository](https://github.com/budimanjojo/talhelper).

- Add the repository as input in your `flake.nix` file

    ```nix
    {
      inputs = {
        talhelper.url = "github:budimanjojo/talhelper";
      }
    }
    ```

- The package is now available at `packages.<system>.default` of the flake. You can call it in your `home.packages` or `environment.systemPackages` or `devShell` by referencing the input as `inputs.talhelper.packages.<system>.default`.

## Using AUR

You can get `talhelper` from [AUR](https://aur.archlinux.org/packages/talhelper-bin) using any [AUR helper](https://wiki.archlinux.org/title/AUR_helpers) if you're Arch Linux user btw.

Example using [`yay`](https://github.com/Jguer/yay):
```bash
yay -S talhelper-bin
```

## Using Scoop

You can get `talhelper` from [Scoop](https://scoop.sh/) if you're a Windows user (thanks to [@dedene](https://github.com/dedene)).

```powershell
scoop bucket add budimanjojo https://github.com/budimanjojo/talhelper.git
scoop install talhelper
```

## Using one liner with jpillora

You can get `talhelper` using this one liner using tool provided by [jpillora](https://github.com/jpillora/installer).

```bash
curl https://i.jpillora.com/budimanjojo/talhelper! | sudo bash
```

## From the release page

If none of the above works for you, you can download the archived binary for your system from the [latest release page](https://github.com/budimanjojo/talhelper/releases/latest).

Please let me know if you want to help with adding new installation method by [creating a new issue](https://github.com/budimanjojo/talhelper/issues/new).
