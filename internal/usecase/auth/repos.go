package auth

import (
	"context"

	domaintoken "github.com/daniyarsan/auth-service/internal/domain/token"
	domainuser "github.com/daniyarsan/auth-service/internal/domain/user"
)

// keep interfaces narrow for testability
type UserRepo interface {
	Create(ctx context.Context, u *domainuser.User) error
	GetByEmail(ctx context.Context, email string) (*domainuser.User, error)
	GetByID(ctx context.Context, id int64) (*domainuser.User, error)
}

type RefreshTokenRepo interface {
	Save(ctx context.Context, rt *domaintoken.RefreshToken) error
	Exists(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, token string) error
	DeleteAllForUser(ctx context.Context, userID int64) error
}
