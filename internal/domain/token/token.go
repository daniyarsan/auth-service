package token

type RefreshToken struct {
	UserID    int64
	Token     string
	ExpiresAt int64 // unix seconds
}
