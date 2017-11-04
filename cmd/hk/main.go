package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Xe/eclier"
	"github.com/Xe/x/tools/glue/libs/gluaexpect"
	"github.com/Xe/x/tools/glue/libs/gluasimplebox"
	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	heroku "github.com/cyberdelia/heroku-go/v3"
	"github.com/dickeyxxx/netrc"
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
	luar "layeh.com/gopher-luar"
)

var (
	scriptHome = flag.String("script-home", "./scripts", "the script home for command handlers")
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

	n, err := netrc.Parse(filepath.Join(os.Getenv("HOME"), ".netrc"))
	if err != nil {
		panic(err)
	}

	L.SetGlobal("netrc", luar.New(L, n))

	user := n.Machine("api.heroku.com").Get("login")
	pass := n.Machine("api.heroku.com").Get("password")

	if user != "" && pass != "" {
		heroku.DefaultTransport.Username = user
		heroku.DefaultTransport.Password = pass
	}

	h := heroku.NewService(heroku.DefaultClient)

	L.SetGlobal("set_heroku_userpass", luar.New(L, func(user, pass string) {
		heroku.DefaultTransport.Username = user
		heroku.DefaultTransport.Password = pass
	}))
	L.SetGlobal("heroku_new_token", luar.New(L, func() (string, error) {
		description := "hk login from " + time.Now().UTC().Format(time.RFC3339)
		expires := 2592000 // 30 days
		opts := heroku.OAuthAuthorizationCreateOpts{
			Description: &description,
			ExpiresIn:   &expires,
		}

		auth, err := h.OAuthAuthorizationCreate(context.Background(), opts)
		if err != nil {
			return "", err
		}

		return auth.AccessToken.Token, nil
	}))
}
