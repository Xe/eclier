package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/Xe/eclier"
	"github.com/Xe/x/tools/glue/libs/gluaexpect"
	"github.com/Xe/x/tools/glue/libs/gluasimplebox"
	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluaquestion"
	"github.com/kohkimakimoto/gluassh"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/otm/gluaflag"
	"github.com/otm/gluash"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	json "layeh.com/gopher-json"
)

var (
	scriptHome = flag.String("script-home", "/keybase/private/xena/scdo", "the script home for command handlers")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	r, err := eclier.NewRouter(
		eclier.WithGluaCreationHook(preload),
		eclier.WithScriptHome(*scriptHome),
	)
	if err != nil {
		panic(err)
	}

	r.Run(ctx, flag.Args())
}

func preload(L *lua.LState) {
	L.PreloadModule("re", gluare.Loader)
	L.PreloadModule("sh", gluash.Loader)
	L.PreloadModule("fs", gluafs.Loader)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("question", gluaquestion.Loader)
	L.PreloadModule("ssh", gluassh.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("flag", gluaflag.Loader)
	L.PreloadModule("template", gluatemplate.Loader)
	L.PreloadModule("url", gluaurl.Loader)
	gluaexpect.Preload(L)
	gluasimplebox.Preload(L)
	gluaxmlpath.Preload(L)
	json.Preload(L)
}
