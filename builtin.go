package eclier

import (
	"context"
	"flag"
	"fmt"
)

// Constants for built-in commands.
const (
	BuiltinScriptPath = "<built-in>"
	BuiltinAuthor     = "<built-in>"
	BuiltinVersion    = "<built-in>"
)

type pluginCommand struct {
	r  *Router
	fs *flag.FlagSet

	dontShowBuiltin *bool
}

// Close is a no-op.
func (p *pluginCommand) Close() error { return nil }

// Init sets up the flags for this command.
func (p *pluginCommand) Init() {
	p.fs = flag.NewFlagSet(p.Verb(), flag.ExitOnError)

	p.dontShowBuiltin = p.fs.Bool("no-builtin", false, "if set, don't show built-in commands")
}

// ScriptPath returns the built-in script path.
func (p *pluginCommand) ScriptPath() string { return BuiltinScriptPath }

// Verb returns the command verb.
func (p *pluginCommand) Verb() string { return "plugins" }

// Help returns the command help
func (p *pluginCommand) Help() string {
	return `plugin lists all of the loaded commands and their script paths.`
}

func (p *pluginCommand) Usage() string {
	return `  -no-builtin 
    	if set, don't show built-in commands`
}

func (p *pluginCommand) Author() string { return BuiltinAuthor }

func (p *pluginCommand) Version() string { return BuiltinVersion }

// Run executes the command.
func (p *pluginCommand) Run(ctx context.Context, arg []string) error {
	p.fs.Parse(arg)

	for _, c := range p.r.cmds {
		if c.ScriptPath() == BuiltinScriptPath && *p.dontShowBuiltin {
			continue
		}

		fmt.Printf("%s\t%s\n", c.Verb(), c.ScriptPath())
	}

	return nil
}

// NewBuiltinCommand makes it easier to write core commands for eclier.
func NewBuiltinCommand(verb, help, usage string, doer func(context.Context, []string) error) Command {
	return &commandFunc{
		verb:  verb,
		help:  help,
		usage: usage,
		doer:  doer,
	}
}

// commandFunc is a simple alias for creating builtin commands.
type commandFunc struct {
	verb  string
	help  string
	usage string
	doer  func(context.Context, []string) error
}

// Close deallocates resources set up by the initialization of the command.
func (c *commandFunc) Close() error { return nil }

// Init is a no-op.
func (c *commandFunc) Init() {}

// ScriptPath returns the built-in script path.
func (c *commandFunc) ScriptPath() string { return BuiltinScriptPath }

// Verb returns the command verb.
func (c *commandFunc) Verb() string { return c.verb }

// Help returns the command help.
func (c *commandFunc) Help() string { return c.help }

// Usage returns the command usage.
func (c *commandFunc) Usage() string { return c.usage }

// Author returns the built-in author.
func (c *commandFunc) Author() string { return BuiltinAuthor }

// Version returns the built-in version.
func (c *commandFunc) Version() string { return BuiltinVersion }

// Run runs the command handler.
func (c *commandFunc) Run(ctx context.Context, arg []string) error {
	return c.doer(ctx, arg)
}
