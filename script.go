package eclier

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"sync"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type Script struct {
	Verb    string
	Help    string
	Usage   string
	Author  string
	Version string
}

type gluaCommand struct {
	sync.Mutex
	script *Script
	L      *lua.LState

	filename string
}

func newGluaCommand(preload func(*lua.LState), filename string, r io.Reader) Command {
	L := lua.NewState()
	preload(L)

	script := &Script{}
	L.SetGlobal("script", luar.New(L, script))

	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	err = L.DoString(string(data))
	if err != nil {
		panic(err)
	}

	return &gluaCommand{script: script, L: L, filename: filename}
}

func (g *gluaCommand) Close() error {
	g.L.Close()
	return nil
}

func (g *gluaCommand) Init() {}

func (g *gluaCommand) ScriptPath() string { return g.filename }

func (g *gluaCommand) Verb() string { return g.script.Verb }

func (g *gluaCommand) Help() string { return g.script.Help }

func (g *gluaCommand) Usage() string { return g.script.Usage }

func (g *gluaCommand) Author() string { return g.script.Author }

func (g *gluaCommand) Version() string { return g.script.Version }

func (g *gluaCommand) Run(ctx context.Context, arg []string) error {
	runf := g.L.GetGlobal("run")

	if runf.Type() == lua.LTNil {
		return errors.New("no global function run in this script")
	}

	return g.L.CallByParam(lua.P{
		Fn:      runf,
		NRet:    0,
		Protect: true,
	}, luar.New(g.L, arg))
}
