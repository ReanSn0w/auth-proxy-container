package server

import (
	"net/http"
	"net/http/httputil"

	"github.com/go-pkgz/lgr"
)

func (s *Server) reverseProxy() http.Handler {
	proxy := &httputil.ReverseProxy{
		ErrorLog: lgr.ToStdLogger(s.log, "WARN"),
		Director: func(r *http.Request) {
			r.URL.Scheme = s.proxy.Scheme
			r.URL.Host = s.proxy.Host
			r.Header.Set("X-User-ID", getUserID(r.Context()))
		},
	}
	return proxy
}
