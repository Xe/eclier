package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/Xe/eclier"
	"github.com/Xe/eclier/internal/gluaflag"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/otm/gluash"
	lua "github.com/yuin/gopher-lua"
	json "layeh.com/gopher-json"
)

var (
	scriptHome = flag.String("script-home", "./scripts", "the script home for command implementations")
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
	L.PreloadModule("sh", gluash.Loader)
	L.PreloadModule("fs", gluafs.Loader)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("flag", gluaflag.Loader)
	L.PreloadModule("template", gluatemplate.Loader)
	L.PreloadModule("url", gluaurl.Loader)
	json.Preload(L)
}
