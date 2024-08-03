package postgres

import (
	"context"
	"database/sql"
	"sso/m/internal/config"
	"sso/m/internal/domain/models"
	"strconv"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func New(cfg config.PostgresConfig) *Postgres {
	db, err := sql.Open("postgres", "postgresql://"+cfg.User+":"+cfg.Password+"@"+cfg.Host+":"+strconv.Itoa(cfg.Port)+"/"+cfg.Database+"?sslmode=disable")

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Close() {
	p.db.Close()
}

func (p *Postgres) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRowContext(ctx, email, passwordHash).Scan(&id)

	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Postgres) FindUser(ctx context.Context, email string) (*models.User, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "SELECT id, email, password_hash FROM users WHERE email = $1")
	if err != nil {
		return nil, err
	}

	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Postgres) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "SELECT is_admin FROM users WHERE id = $1")
	if err != nil {
		return false, err
	}

	var isAdmin bool
	err = stmt.QueryRowContext(ctx, userId).Scan(&isAdmin)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return isAdmin, nil
}

func (p *Postgres) FindApp(ctx context.Context, appID int) (*models.App, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return nil, err
	}

	var app models.App
	err = stmt.QueryRowContext(ctx, appID).Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &app, nil
}