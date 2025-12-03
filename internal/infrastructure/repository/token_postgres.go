package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/daniyarsan/auth-service/internal/domain/token"
)

type TokenPostgres struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *TokenPostgres {
	return &TokenPostgres{db: db}
}

func (r *TokenPostgres) Save(ctx context.Context, rt *token.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(ctx, query, rt.Token, rt.UserID, rt.ExpiresAt)
	return err
}

func (r *TokenPostgres) Exists(ctx context.Context, tokenStr string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM refresh_tokens WHERE token = $1)`
	if err := r.db.QueryRow(ctx, query, tokenStr).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *TokenPostgres) Delete(ctx context.Context, tokenStr string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(ctx, query, tokenStr)
	return err
}

func (r *TokenPostgres) DeleteAllForUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
