package eclier

import (
	"context"
)

// Command is an individual subcommand.
type Command interface {
	Close() error
	Init()
	ScriptPath() string
	Verb() string
	Help() string
	Run(ctx context.Context, arg []string) error
}
