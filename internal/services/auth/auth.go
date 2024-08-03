package auth

import (
	"context"
	"log/slog"
	"sso/m/internal/domain/models"
	"sso/m/internal/lib/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	UserProvider
	UserSaver
	AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error)
}

type UserProvider interface {
	FindUser(ctx context.Context, email string) (*models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	FindApp(ctx context.Context, appID int) (*models.App, error)
}

func NewAuth(
	log *slog.Logger,
	userProvider UserProvider,
	userSaver UserSaver,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		UserProvider: userProvider,
		UserSaver:    userSaver,
		AppProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (app *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {

	user, err := app.UserProvider.FindUser(ctx, email)
	if err != nil {
		app.log.Error("failed to find user", slog.Any("err", err.Error()))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		app.log.Error("failed to compare password", slog.Any("err", err.Error()))
		return "", err
	}

	application, err := app.AppProvider.FindApp(ctx, appID)
	if err != nil {
		app.log.Error("failed to find app", slog.Any("err", err.Error()))
		return "", err
	}

	token, err := jwt.NewToken(user, application, app.tokenTTL)
	if err != nil {
		app.log.Error("failed to create token", slog.Any("err", err.Error()))
		return "", err
	}

	return token, nil
}

func (app *Auth) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.log.Error("failed to hash password", slog.Any("err", err.Error()))
		return 0, err
	}

	userID, err := app.UserSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		app.log.Error("failed to save user", slog.Any("err", err.Error()))
		return 0, err
	}

	return userID, nil
}

func (app *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	isAdmin, err := app.UserProvider.IsAdmin(ctx, userID)
	if err != nil {
		app.log.Error("failed to find user", slog.Any("err", err.Error()))
		return false, err
	}

	return isAdmin, nil
}