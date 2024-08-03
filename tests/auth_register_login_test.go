package tests

import (
	"sso/m/tests/suite"
	"testing"
	"time"

	ssov1 "github.com/EthernetUser/sso-protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID        = 0
	appId             = 1
	appSecret         = "secret"
	passDefaultLength = 10
)

func TestRegisterLogin(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()
	respReg, err := st.AuthClient.Register(ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: password,
		},
	)

	require.NoError(t, err)

	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx,
		&ssov1.LoginRequest{
			Email:    email,
			Password: password,
			AppId:    appId,
		},
	)
	require.NoError(t, err)
	token := respLogin.GetToken()
	assert.NotEmpty(t, token)
	loginTime := time.Now()

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})

	require.NoError(t, err)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["userEmail"].(string))
	assert.Equal(t, appId, int(claims["appId"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["userId"].(float64)))

	deltaSeconds := 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), float64(deltaSeconds))
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLength)
}
