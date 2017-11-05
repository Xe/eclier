package gluanetrc

import (
	"os"
	"path/filepath"

	"github.com/dickeyxxx/netrc"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

var n *netrc.Netrc

func init() {
	var err error

	n, err = netrc.Parse(filepath.Join(os.Getenv("HOME"), ".netrc"))
	if err != nil {
		panic(err)
	}
}

var exports = map[string]lua.LGFunction{
	"machine":       machine,
	"save":          save,
	"removeMachine": removeMachine,
}

func removeMachine(L *lua.LState) int {
	name := L.ToString(1)

	n.RemoveMachine(name)

	return 0
}

func machine(L *lua.LState) int {
	name := L.ToString(1)

	m := n.Machine(string(name))

	L.Push(luar.New(L, m))
	return 1
}

func save(L *lua.LState) int {
	n.Save()
	return 0
}

func Preload(L *lua.LState) {
	L.PreloadModule("netrc", Loader)
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
