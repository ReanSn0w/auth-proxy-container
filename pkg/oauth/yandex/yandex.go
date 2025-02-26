package yandex

import (
	"context"
	"net/http"
	"time"

	"github.com/ReanSn0w/auth-proxy-container/pkg/oauth"
	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

func New(clientID string, clientSecret string) oauth.Provider {
	return &YandexClient{
		base: oauth.New(
			clientID, clientSecret,
			"https://oauth.yandex.ru/authorize",
			"https://oauth.yandex.ru/device/code",
			"https://oauth.yandex.ru/token",
			oauth2.AuthStyleAutoDetect),

		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

type YandexClient struct {
	client *http.Client
	base   *oauth.Base
}

func (yc *YandexClient) Link(state, redirectURI string) string {
	return yc.base.Link(state, redirectURI)
}

func (yc *YandexClient) Validate(token string) error {
	return yc.base.Validate(token)
}

func (yc *YandexClient) ExtractUserID(token string) (string, error) {
	return yc.base.ExtractUserID(token)
}

func (yc *YandexClient) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := yc.base.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return yc.getUserJWT(token)

}

func (yc *YandexClient) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	token, err := yc.base.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return yc.getUserJWT(token)
}

type UserInfo struct {
	Login    string `json:"login"`
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Psuid    string `json:"psuid"`
}

func (yc *YandexClient) getUserJWT(token *oauth2.Token) (*oauth2.Token, error) {
	resp := UserInfo{}
	err := web.NewJsonRequest(yc.client, "https://login.yandex.ru/info").
		SetHeader("Authorization", "Oauth "+token.AccessToken).
		SetQuery("format", "json").
		Do(&resp)

	if err != nil {
		return nil, err
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims = jwt.MapClaims{
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour).Unix(),
		"sub":       resp.Login,
		"id":        resp.ID,
		"client_id": resp.ClientID,
		"psuid":     resp.Psuid,
	}

	resignedAccessToken, err := jwtToken.SignedString([]byte(yc.base.Config.ClientSecret))
	if err != nil {
		return nil, err
	}

	token.AccessToken = resignedAccessToken
	return token, nil
}
