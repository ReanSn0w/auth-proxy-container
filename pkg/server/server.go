package server

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/ReanSn0w/auth-proxy-container/pkg/oauth"
	"github.com/ReanSn0w/auth-proxy-container/pkg/utils"
	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/go-chi/chi"
	"github.com/go-pkgz/lgr"
)

func New(log lgr.L, pathRE, uidRE *regexp.Regexp, proxy *url.URL, provider oauth.Provider) *Server {
	return &Server{
		log:      log,
		provider: provider,
		proxy:    proxy,
		pathRE:   pathRE,
		uidRE:    uidRE,
		state:    utils.New(),
		srv:      web.New(log),
	}
}

type Server struct {
	log   lgr.L
	state *utils.Storage
	srv   *web.Server

	provider oauth.Provider
	proxy    *url.URL
	pathRE   *regexp.Regexp
	uidRE    *regexp.Regexp
}

func (s *Server) Start(port int) {
	cancel := func(err error) {
		if err != nil {
			s.log.Logf("[ERROR] server start err: %v", err)
			os.Exit(2)
		}
	}

	s.srv.Run(cancel, port, s.handler())
}

func (s *Server) Stop(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Logf("[ERROR] shutting down server: %v", err)
	}
}

func (s *Server) handler() http.Handler {
	r := chi.NewRouter()

	r.Use(s.checkToken)

	r.Route("/oauth", func(r chi.Router) {
		r.Get("/authorize", s.authorize)
	})

	r.Handle("/*", s.reverseProxy())

	return r
}
