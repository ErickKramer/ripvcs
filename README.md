# 🎒 ripvcs

**PRE-RELEASE: Heavily developing this tool at the moment.**

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

1. Go to the [Releases](https://github.com/ErickKramer/ripvcs/releases) page.
2. Download latest binary
3. Make the binary executable
   ```console
   chmod +x rv
   ```
4. Move binary to a directory included in your `PATH`
   ```console
   mv rv ~/.local/bin/
   ```
5. Verify installation
   ```console
   rv help
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
3. Follow steps 3-5 from the "Using pre-built binaries" section.

## Usage

To check the available commands in _rv_ simply run `rv help`:

```console
rv help

Fast CLI tool for managing multiple Git repositories.

Usage:
  rv [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
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
  -h, --help                  help for import
  -i, --input .repos          Path to input .repos file
  -r, --recursive .repos      Recursively search of other .repos file in the cloned repositories
  -l, --shallow               Clone repositories with a depth of 1
  -s, --skip                  Skip existing repositories
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
```

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

|    Command    |   vcs   |   rv    |
| :-----------: | :-----: | :-----: |
|    import     | 2.363 s | 1.753 s |
| import + skip | 0.691 s | 0.004 s |
|      log      | 0.248 s | 0.021 s |
|     pull      | 0.635 s | 0.417 s |
|    status     | 0.238 s | 0.035 s |
|   validate    | 0.869 s | 0.414 s |

## Future enhancements

- [ ] Support tar artifacts
- [ ] Support export command
- [ ] Support custom git commands
