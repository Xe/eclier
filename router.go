package eclier

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

// Router is the main subcommand router for eclier. At a high level what this is
// doing is similar to http.ServeMux, but for CLI commands instead of HTTP handlers.
type Router struct {
	lock sync.Mutex
	cmds map[string]Command

	// configured data
	gluaCreationHook func(*lua.LState)
	scriptHome       string
	cartridge        map[string]string
}

// NewRouter creates a new instance of Router and sets it up for use.
func NewRouter(opts ...RouterOption) (*Router, error) {
	r := &Router{
		cmds:      map[string]Command{},
		cartridge: map[string]string{},
	}

	for _, opt := range opts {
		opt(r)
	}

	// scan r.scriptHome for lua scripts, load them into their own lua states and
	// make a wrapper around them for the Command type.

	err := filepath.Walk(r.scriptHome, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("error in arg: %v", err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".lua") {
			c := newGluaCommand(r.gluaCreationHook, filepath.Join(r.scriptHome, info.Name()))

			r.cmds[c.Verb()] = c
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	r.cmds["plugin"] = &pluginCommand{r: r}

	return r, nil
}

// Run executes a single command given in slot 0 of the argument array.
func (r *Router) Run(ctx context.Context, arg []string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(arg) == 0 {
		fmt.Printf("please specify a subcommand, such as `%s help`\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	cmd := arg[0]
	arg = arg[1:]

	ci, ok := r.cmds[cmd]
	if !ok {
		fmt.Printf("No such command %s could be run.\n", cmd)
		os.Exit(2)
	}

	ci.Init()
	return ci.Run(ctx, arg)
}

// RouterOption is a functional option for Router.
type RouterOption func(*Router)

// WithScriptHome sets the router's script home to the given directory. This is
// where lua files will be walked and parsed.
func WithScriptHome(dir string) RouterOption {
	return func(r *Router) {
		r.scriptHome = dir
	}
}

// WithGluaCreationHook adds a custom bit of code that runs every time a new
// gopher-lua LState is created. This allows users of this library to register
// custom libraries to the pile of states.
func WithGluaCreationHook(hook func(*lua.LState)) RouterOption {
	return func(r *Router) {
		r.gluaCreationHook = hook
	}
}
