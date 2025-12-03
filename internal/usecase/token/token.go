package token

import "context"

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenManager interface {
	Generate(ctx context.Context, userID int64) (Tokens, error)
	ParseAccessToken(ctx context.Context, token string) (int64, error) // returns userID
}
