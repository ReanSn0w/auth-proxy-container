package server

import (
	"errors"
	"net/http"

	"github.com/ReanSn0w/gokit/pkg/web"
	"golang.org/x/oauth2"
)

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	token, err := s.provider.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeTokens(w, token, r.URL.Scheme == "https")

	url, ok := s.state.Fire(state)
	if !ok {
		web.NewResponse(errors.New("invalid state"))
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) writeTokens(w http.ResponseWriter, token *oauth2.Token, isSecureHost bool) {
	accessToken := &http.Cookie{
		Name:     accessToken,
		Value:    token.AccessToken,
		HttpOnly: true,
		Secure:   isSecureHost,
		Path:     "/",
	}

	http.SetCookie(w, accessToken)

	refreshToken := &http.Cookie{
		Name:     refreshToken,
		Value:    token.RefreshToken,
		HttpOnly: true,
		Secure:   isSecureHost,
		Path:     "/",
	}

	http.SetCookie(w, refreshToken)
}
