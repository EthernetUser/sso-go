package jwt

import (
	"sso/m/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user *models.User, app *models.App, tokenTTL time.Duration) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(tokenTTL).Unix()
	claims["appId"] = app.Id
	claims["userId"] = user.Id
	claims["userEmail"] = user.Email

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}