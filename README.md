# `docsync`

Documentation synchronization tool and library. 

Tool is wrote using [cobra](https://github.com/spf13/cobra) and [promptui](https://github.com/manifoldco/promptui), so it works like command line command and can run interactively.

```shell
❯ ./docsync --help                                               
docsync is a documentation synchronization tool.

docsync creates a template for the new translation of documentation to your project into another language.

Usage:
  docsync [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create new documentation from base lang. Command should be called from the root documentation directory.
  help        Help about any command
  init        Initialize new documentation.
  sync        Synchronizes specified language documentation with base. Command should be called from the root documentation directory.
  update      Update hashkeys of documentation blocks in base documentation. Command should be called from the root documentation directory.

Flags:
  -h, --help          help for docsync
      --no-interact   Ask for missing flags in interactive mod or not.
  -v, --version       version for docsync

Use "docsync [command] --help" for more information about a command.
```

## Usage

Firstly initialize a new documentation.

```shell
❯ ./docsync init --help
Initialize new documentation.

Usage:
  docsync init [flags]

Flags:
  -h, --help             help for init
      --lang string      The base language of the documentation. Please, specify ISO639-1 code. (default "en")
      --path string      Path where docs will be placed. '.' means the current directory. (default ".")
      --plangs strings   Programming languages of code from the documentations.
      --type string      The main documentation files' type. (default "md")

Global Flags:
      --no-interact   Ask for missing flags in interactive mod or not.
```

You will obtain a documentation template in `/path/<base>/index.<type>` (for example `/docs/en/index.md`). The `path/docsync.yaml` is the `docsync` configuration file.

After you wrote documentation, you should run `docsync update` from the directory, where `docsync.yaml` is located. This will update documentation blocks' identifiers.

```shell
❯ ./docsync update --help
Update hashkeys of documentation blocks in base documentation. Command should be called from the root documentation directory.

Usage:
  docsync update [flags]

Flags:
  -h, --help            help for update
      --update-status   Update status of 'New' documentation blocks.

Global Flags:
      --no-interact   Ask for missing flags in interactive mod or not.
```

Now, the documentation in another language can be created from the base documentation (or empty documentation).

```shell
❯ ./docsync create --help
Create new documentation from base lang. Command should be called from the root documentation directory.

Usage:
  docsync create [flags]

Flags:
      --create-empty       Create an empty documentation.
      --create-from-base   Create documentation using base as template. (default true)
  -h, --help               help for create
  -l, --lang string        New documentation language.

Global Flags:
      --no-interact   Ask for missing flags in interactive mod or not.
```

If base documentation is updated, you may wish to update the documentation in another language without eliminating existing one.

```shell
❯ ./docsync sync --help  
Synchronizes specified language documentation with base. Command should be called from the root documentation directory.

Usage:
  docsync sync [flags]

Flags:
  -h, --help          help for sync
  -l, --lang string   Documentation language for synchronization.

Global Flags:
      --no-interact   Ask for missing flags in interactive mod or not.
```

## Future work

For now the tool supports only documentation in markdown files. I plan to add support for .rst and .xml documentations.

## License

This project is licensed under the [MIT license](https://github.com/Scurrra/docsync/blob/master/LICENSE).

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in docsync by you, shall be licensed as MIT, without any additional terms or conditions. 

If you have any issues or you find any mistakes or bugs, contact me. 