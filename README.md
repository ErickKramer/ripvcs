# 🎒 ripvcs

<p align="center">
  <img width="33%" src="./assets/cat_sorter.jpeg" />
</p>

<!-- prettier-ignore -->
_ripvcs (rv)_ is a command-line tool written in Go, providing an efficient alternative to [vcstool](https://github.com/dirk-thomas/vcstool) for managing multiple repository workspaces.

Whether you are managing a few repositories or a complex workspace with numerous nested repositories, ripvcs (rv) offers the performance and efficiency you need to keep your workflow smooth and responsive.

## 🪄 Features

- **Enhanced Concurrency:** Utilizes Go routines to manage multiple tasks simultaneously,
  ensuring optimal performance and reducing wait times.
- **Recursive Import:** Supports recursive import functionality to automatically search for `.repos`
  files within directories to streamline the process of managed nested
  repositories dependencies.

## 🧰 Installation

### Using pre-built binaries

- Latest release <https://github.com/ErickKramer/ripvcs/releases/latest>

```console
RIPVCS_VERSION=$(curl -s "https://api.github.com/repos/ErickKramer/ripvcs/releases/latest" | \grep -Po '"tag_name": *"v\K[^"]*')
ARCHITECTURE="linux_amd64"
curl -Lo ~/.local/bin/rv "https://github.com/ErickKramer/ripvcs/releases/download/v${RIPVCS_VERSION}/ripvcs_${RIPVCS_VERSION}_${ARCHITECTURE}"
chmod +x ~/.local/bin/rv
```

## Build from source

1. Clone repository
   ```console
   git clone https://github.com/ErickKramer/ripvcs
   cd ripvcs
   ```
2. Build binary
   ```console
   go build -o rv main.go
   ```
3. Move rv to path
   ```console
   mv rv ~/.local/bin/rv
   chmod +x ~/.local/bin/rv
   ```
## Usage

To check the available commands in _rv_ simply run `rv help`:

```console
rv help

Fast CLI tool for managing multiple Git repositories.

Usage:
  rv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  export      Export list of available repositories
  help        Help about any command
  import      Import repositories listed in the given .repos file
  log         Get logs of all repositories.
  pull        Pull latest version from remote.
  status      Check status of all repositories
  switch      Switch repository version
  sync        Synchronize all found repositories.
  validate    Validate a .repos file
  version     Print the version number
```

Each of the available commands have their own help with information about their usage and available flags (e.g. `rv help import`).

```console
Import repositories listed in the given .repos file

The repositories are cloned in the given path or in the current path.

It supports recursively searching for any other .repos file found at each
import cycle.

Usage:
  rv import <optional path> [flags]

Flags:
  -d, --depth-recursive int   Regulates how many levels the recursive dependencies would be cloned. (default -1)
  -x, --exclude strings       List of files and/or directories to exclude when performing a recursive import
  -f, --force                 Force overwriting existing repositories
  -h, --help                  help for import
  -i, --input .repos          Path to input .repos file
  -s, --recurse-submodules    Recursively clone submodules
  -r, --recursive .repos      Recursively search of other .repos file in the cloned repositories
  -n, --retry int             Number of attempts to import repositories (default 2)
  -l, --shallow               Clone repositories with a depth of 1
  -w, --workers int           Number of concurrent workers to use (default 8)
```

### Repositories file

```yaml
repositories:
  demos:
    type: git
    url: https://github.com/ros2/demos
    version: jazzy
  stable_demos:
    type: git
    url: https://github.com/ros2/demos
    version: 0.20.4
  default_demos:
    type: git
    url: https://github.com/ros2/demos
    exclude: []
```

### Import exclusion

It is possible to exclude files or directories when doing recursive import. This can be done either
through the use of the `--exclude / -x` flag accompanied by the name of the `.repos` file or a
directory within the path.

Additionally, it is possible to add a `exclude` attribute to the `.repos` file to hard-code what
files to exclude during import. An example of this can be seen in [nested_example.repos](./test/nested_example.repos)

## Related Project

- [vcstool](https://github.com/dirk-thomas/vcstool)

## Shell completions

`rv` supports generating an autocompletion script for bash, fish, powershell, and zsh

For example, to configure to generate the completion for zsh do the following

```console
rv completion zsh _ripvcs
```

Then, you need to place the completion file in the proper location to be loaded by your zsh configuration.

Afterwards, you should be able to do `rv <TAB>` to get autocompletion for the available commands and `rv import -<TAB>` to get autocompletion for the available flags.

### Runtime comparison

I ran a quick comparison between `rv` and `vcs` using the [valid_example.repos](./test/valid_example.repos) file for different commands using `8` workers.

|      Command       |   vcs   |   rv    |
| :----------------: | :-----: | :-----: |
| import (overwrite) | 2.363 s | 1.753 s |
|   import + skip    | 0.691 s | 0.004 s |
|        log         | 0.248 s | 0.021 s |
|        pull        | 0.635 s | 0.417 s |
|       status       | 0.238 s | 0.035 s |
|      validate      | 0.869 s | 0.414 s |

## Future enhancements

- [ ] Support tar artifacts
- [ ] Support export command
- [ ] Support custom git commands
