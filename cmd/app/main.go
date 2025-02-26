package main

import (
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/ReanSn0w/auth-proxy-container/pkg/oauth"
	"github.com/ReanSn0w/auth-proxy-container/pkg/oauth/yandex"
	"github.com/ReanSn0w/auth-proxy-container/pkg/server"
	"github.com/ReanSn0w/gokit/pkg/app"
)

var (
	revision = "unknown"
	opts     = struct {
		app.Debug

		Port    int    `short:"p" long:"port" env:"PORT" default:"8080" description:"application listen port"`
		Output  string `short:"o" long:"output" env:"OUTPUT" default:"https://yandex.ru" description:"output URL"`
		Private string `long:"private" env:"PRIVATE" default:"^/(.*)" description:"private path regexp"`
		UserID  string `long:"uid" env:"USER_ID" default:"none" description:"user ID regexp"`

		Oauth struct {
			Provider     string `long:"provider" env:"PROVIDER" default:"yandex" choice:"yandex" description:"OAuth provider"`
			ClientID     string `long:"client-id" env:"CLIENT_ID" description:"OAuth client ID"`
			ClientSecret string `long:"client-secret" env:"CLIENT_SECRET" description:"OAuth client secret"`
		} `group:"oauth" namespace:"oauth" env-namespace:"OAUTH"`
	}{}
)

func main() {
	app := app.New("Auth Proxy Container", revision, &opts)

	{
		privatePath, err := regexp.Compile(opts.Private)
		if err != nil {
			app.Log().Logf("[ERROR] compile private path regexp error: %v", err)
			os.Exit(2)
			return
		}

		privateUID, err := regexp.Compile(opts.UserID)
		if err != nil {
			app.Log().Logf("[ERROR] compile user ID regexp error: %v", err)
			os.Exit(2)
			return
		}

		proxyURL, err := url.Parse(opts.Output)
		if err != nil {
			app.Log().Logf("[ERROR] parse output URL error: %v", err)
			os.Exit(2)
			return
		}

		if opts.Oauth.ClientID == "" || opts.Oauth.ClientSecret == "" {
			app.Log().Logf("[ERROR] OAuth client ID or secret is empty")
			os.Exit(2)
			return
		}

		var provider oauth.Provider
		switch opts.Oauth.Provider {
		case "yandex":
			provider = yandex.New(opts.Oauth.ClientID, opts.Oauth.ClientSecret)
		default:
			app.Log().Logf("[ERROR] unknown OAuth provider: %s", opts.Oauth.Provider)
			os.Exit(2)
			return
		}

		srv := server.New(app.Log(), privatePath, privateUID, proxyURL, provider)
		app.Add(srv.Stop)
		srv.Start(opts.Port)
	}

	app.GracefulShutdown(time.Second * 10)
}
