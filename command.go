package eclier

import (
	"context"
	"io"
)

// Command is an individual subcommand.
type Command interface {
	io.Closer

	Init()
	ScriptPath() string
	Verb() string
	Help() string
	Run(ctx context.Context, arg []string) error
}
