package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

type Provider interface {
	Link(state, redirectURI string) string
	Validate(string) error
	ExtractUserID(token string) (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error)
}
