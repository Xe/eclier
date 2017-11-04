# eclier

The core of a command line application allowing for trivial user extension.

Every command and subcommand is its own `.lua` file that is either shipped as
part of the built-in cartridge of commands or a plugin that the user installs.

The core contains the following:

- A module loading system for preparing different commands for use
- A stateful layer for plugins to persist data
- The core subcommand router
