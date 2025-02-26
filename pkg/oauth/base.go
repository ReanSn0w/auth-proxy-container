package oauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

func New(clientID, clientSecret, authURL, deviceAuthURL, tokenURL string, authStyle oauth2.AuthStyle) *Base {
	return &Base{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:       authURL,
				TokenURL:      tokenURL,
				DeviceAuthURL: deviceAuthURL,
				AuthStyle:     authStyle,
			},
		},
	}
}

type Base struct {
	Config *oauth2.Config
}

// Создает ссылку для авторизации пользователя на сайте
func (b *Base) Link(state, redirectURI string) string {
	return b.Config.AuthCodeURL(state, oauth2.SetAuthURLParam("redirect_uri", redirectURI))
}

// Производит обмен кода на токен авторизации
func (b *Base) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return b.Config.Exchange(ctx, code)
}

// Производит обмен токена на новый токен авторизации
func (b *Base) Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	return b.Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
}

// Производит проверку токена авторизации на его валидность
// Возвращает ошибку, если токен не валиден
func (b *Base) Validate(token string) error {
	res, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(b.Config.ClientSecret), nil
	})
	if err != nil {
		return err
	}
	if !res.Valid {
		return errors.New("invalid token")
	}
	return nil
}

// Извлекает идентификатор пользователя из токена авторизации
func (b *Base) ExtractUserID(token string) (string, error) {
	res, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(b.Config.ClientSecret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := res.Claims.(jwt.MapClaims); ok && res.Valid {
		if userID, ok := claims["sub"].(string); ok {
			return userID, nil
		}
	}

	return "", fmt.Errorf("invalid token")
}
