package token

import "context"

type Repository interface {
	Save(ctx context.Context, rt *RefreshToken) error
	Exists(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string) error
	DeleteAllForUser(ctx context.Context, userID int64) error
}
