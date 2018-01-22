package eclier

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
	lua "github.com/yuin/gopher-lua"
	"layeh.com/asar"
)

// Router is the main subcommand router for eclier. At a high level what this is
// doing is similar to http.ServeMux, but for CLI commands instead of HTTP handlers.
type Router struct {
	lock sync.Mutex
	cmds map[string]Command

	// configured data
	gluaCreationHook func(*lua.LState)
	scriptHomes      []string
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

	for _, home := range r.scriptHomes {
		err := filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("error in arg: %v", err)
				return err
			}

			if strings.HasSuffix(info.Name(), ".lua") {
				fname := filepath.Join(home, info.Name())
				fin, err := os.Open(fname)
				if err != nil {
					return err
				}
				defer fin.Close()

				c := newGluaCommand(r.gluaCreationHook, fname, fin)
				r.cmds[c.Verb()] = c
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	var helpCommand Command = NewBuiltinCommand("help", "shows help for subcommands", "help [subcommand]", func(ctx context.Context, arg []string) error {
		if len(arg) == 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Verb", "Author", "Version", "Help"})

			for _, cmd := range r.cmds {
				table.Append([]string{cmd.Verb(), cmd.Author(), cmd.Version(), cmd.Help()})
			}

			table.Render()
			return nil
		}

		cmd, ok := r.cmds[arg[0]]
		if !ok {
			fmt.Printf("can't find help for %s", arg[0])
			os.Exit(2)
		}

		fmt.Printf("Verb: %s\nAuthor: %s\nVersion: %s\nHelp: %s\nUsage: %s %s\n", cmd.Verb(), cmd.Author(), cmd.Version(), cmd.Help(), cmd.Verb(), cmd.Usage())
		return nil
	})

	r.cmds["plugins"] = &pluginCommand{r: r}
	r.cmds["help"] = helpCommand

	return r, nil
}

// AddCommand adds a given command instance to the eclier router.
func (r *Router) AddCommand(cmd Command) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.cmds[cmd.Verb()] = cmd
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
		r.scriptHomes = append(r.scriptHomes, dir)
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

// WithFilesystem loads a http.FileSystem full of lua scripts into this eclier
// router.
func WithFilesystem(shortName string, fs http.FileSystem) RouterOption {
	return func(r *Router) {
		fin, err := fs.Open("/")
		if err != nil {
			log.Fatal(err)
		}
		defer fin.Close()

		childs, err := fin.Readdir(-1)
		if err != nil {
			log.Fatal(err)
		}

		for _, chl := range childs {
			if strings.HasSuffix(chl.Name(), ".lua") {
				fname := filepath.Join(shortName, chl.Name())
				sFin, err := fs.Open(chl.Name())
				if err != nil {
					log.Fatal(err)
				}
				defer sFin.Close()

				c := newGluaCommand(r.gluaCreationHook, fname, sFin)
				r.cmds[c.Verb()] = c
			}
		}
	}
}

// WithAsarFile loads an asar file full of lua scripts into this eclier router.
func WithAsarFile(shortName, fname string) RouterOption {
	return func(r *Router) {
		fin, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer fin.Close()

		e, err := asar.Decode(fin)
		if err != nil {
			log.Fatal(err)
		}

		err = e.Walk(func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".lua") {
				fname := filepath.Join(shortName, "::", path)
				fin := e.Find(path)
				if fin == nil {
					return nil
				}

				c := newGluaCommand(r.gluaCreationHook, fname, fin.Open())
				r.cmds[c.Verb()] = c
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
