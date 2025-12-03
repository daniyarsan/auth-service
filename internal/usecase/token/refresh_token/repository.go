package refresh_token

import (
	"time"

	"gorm.io/gorm"
)

type RefreshTokenEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null;type:text"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Save(token string, userID int64, expiresAt time.Time) error {
	t := &RefreshTokenEntity{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return r.db.Create(t).Error
}

func (r *TokenRepo) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&RefreshTokenEntity{}).Error
}

func (r *TokenRepo) Find(token string) (*RefreshTokenEntity, error) {
	var t RefreshTokenEntity
	if err := r.db.Where("token = ?", token).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
