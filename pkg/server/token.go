package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ReanSn0w/gokit/pkg/web"
)

const (
	accessToken  = "access_token"
	refreshToken = "refresh_token"
)

func (s *Server) checkToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.pathRE.Match([]byte(r.URL.Path)) || strings.HasPrefix(r.URL.Path, "/oauth/") {
			h.ServeHTTP(w, r)
			return
		}

		accessToken, err := r.Cookie(accessToken)
		if err == nil {
			if s.provider.Validate(accessToken.Value) != nil {
				err = errors.New("invalid access token")
			}
		}
		if err != nil {
			refreshToken, err := r.Cookie(refreshToken)
			if err != nil {
				state := s.state.New(r.URL.String())
				url := s.provider.Link(state, generateRedirectURI(r.URL))
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}

			token, err := s.provider.Refresh(r.Context(), refreshToken.Value)
			if err != nil {
				state := s.state.New(r.URL.String())
				url := s.provider.Link(state, generateRedirectURI(r.URL))
				http.Redirect(w, r, url, http.StatusTemporaryRedirect)
				return
			}

			s.writeTokens(w, token, r.URL.Scheme == "https")
			http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
			return
		}

		userID, _ := s.provider.ExtractUserID(accessToken.Value)
		if !s.uidRE.Match([]byte(userID)) {
			web.NewResponse(fmt.Errorf("invalid user id. current: %v", userID)).
				Write(http.StatusUnauthorized, w)
			return
		}

		ctx := setUserID(r.Context(), userID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRedirectURI(url *url.URL) string {
	host := url.Host
	if host == "" {
		host = "localhost"
	}
	scheme := url.Scheme
	if scheme == "" {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%v/oauth/authorize", scheme, host)
}
