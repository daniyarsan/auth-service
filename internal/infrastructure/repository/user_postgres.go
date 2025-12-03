package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/daniyarsan/auth-service/internal/domain/user"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, u *user.User) error {
	query := `
        INSERT INTO users (email, password_hash, created_at, updated_at)
        VALUES ($1,$2,$3,$4)
        RETURNING id
    `
	return r.db.QueryRow(ctx, query,
		u.Email,
		u.PasswordHash,
		u.CreatedAt,
		u.UpdatedAt,
	).Scan(&u.ID)
}

func (r *UserPostgres) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	query := `
        SELECT id, email, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgres) GetByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	query := `
        SELECT id, email, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
    `
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
