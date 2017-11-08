package gluaheroku

import (
	"context"
	"time"

	heroku "github.com/cyberdelia/heroku-go/v3"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

const (
	gitURLSuf = ".git"
)

var s *heroku.Service

func init() {
	s = heroku.NewService(heroku.DefaultClient)
}

var exports = map[string]lua.LGFunction{
	"get_userpass": getUserPass,
	"set_userpass": setUserPass,
	"create_token": createToken,
	"app_create":   appCreate,
	"app_destroy":  appDestroy,
	"app_info":     appInfo,
	"detect_app":   detectApp,
}

func detectApp(L *lua.LState) int {
	appName, err := AppFromGitRemote("heroku")
	if err != nil {
		L.Error(luar.New(L, err.Error()), 1)
		return 0
	}
	L.Push(lua.LString(appName))
	return 1
}

func getUserPass(L *lua.LState) int {
	L.Push(lua.LString(heroku.DefaultTransport.Username))
	L.Push(lua.LString(heroku.DefaultTransport.Password))
	return 2
}

func appInfo(L *lua.LState) int {
	appName := L.ToString(1)

	app, err := s.AppInfo(context.Background(), appName)
	if err != nil {
		L.Error(luar.New(L, err.Error()), 1)
		return 0
	}

	L.Push(luar.New(L, app))
	return 1
}

func appCreate(L *lua.LState) int {
	appName := L.ToString(1)

	var o heroku.AppCreateOpts

	if appName != "" {
		o.Name = heroku.String(appName)
	}

	app, err := s.AppCreate(context.Background(), o)
	if err != nil {
		L.Error(luar.New(L, err.Error()), 1)
		return 0
	}

	L.Push(luar.New(L, app))
	return 1
}

func appDestroy(L *lua.LState) int {
	appName := L.ToString(1)

	app, err := s.AppDelete(context.Background(), appName)
	if err != nil {
		L.Error(luar.New(L, err.Error()), 1)
		return 0
	}

	L.Push(luar.New(L, app))
	return 1
}

func createToken(L *lua.LState) int {
	description := "hk login from " + time.Now().UTC().Format(time.RFC3339)
	expires := 2592000 // 30 days
	opts := heroku.OAuthAuthorizationCreateOpts{
		Description: &description,
		ExpiresIn:   &expires,
	}

	auth, err := s.OAuthAuthorizationCreate(context.Background(), opts)
	if err != nil {
		L.Error(luar.New(L, err.Error()), 1)
		return 0
	}

	L.Push(lua.LString(auth.AccessToken.Token))

	return 1
}

func setUserPass(L *lua.LState) int {
	user := L.ToString(1)
	pass := L.ToString(2)

	heroku.DefaultTransport.Username = user
	heroku.DefaultTransport.Password = pass

	return 0
}

func Preload(L *lua.LState) {
	L.PreloadModule("heroku", Loader)
}

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}
