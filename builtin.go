package eclier

import (
	"context"
	"flag"
	"fmt"
)

// Constants for built-in commands.
const (
	BuiltinScriptPath = "<built-in>"
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
func (p *pluginCommand) Verb() string { return "plugin" }

// Help returns the command help
func (p *pluginCommand) Help() string {
	return `plugin lists all of the loaded commands and their script paths.`
}

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
