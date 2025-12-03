package refresh_token

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	ID           int64  `json:"id"`
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
