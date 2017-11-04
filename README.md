# eclier

The core of a command line application allowing for trivial user extension.

Every command and subcommand is its own `.lua` file that is either shipped as
part of the built-in cartridge of commands or a plugin that the user installs.

The core contains the following:

- A module loading system for preparing different commands for use
- The core subcommand router

## How to write plugins

Create a new file in the script home named after the plugin subcommand, for
example: `scripts/hello.lua`:

```lua
script.verb = "hello"
script.help = "prints everyone's favorite hello world message"
script.author = "Xe" -- put your github username here
script.version = "1.0"
script.usage = ""

function(run) 
  print "Hello, world!"
end
```

And then run it using the example shell cli:

```console
~/go/src/github.com/Xe/eclier:master λ go run ./example/main.go hello
Hello, world!
```

