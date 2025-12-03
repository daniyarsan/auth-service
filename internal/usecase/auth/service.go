package auth

import (
	"context"
	"errors"
	"time"

	domaintoken "github.com/daniyarsan/auth-service/internal/domain/token"
	domainuser "github.com/daniyarsan/auth-service/internal/domain/user"
	jwtinfra "github.com/daniyarsan/auth-service/internal/infrastructure/jwt"
	"github.com/daniyarsan/auth-service/pkg/crypto"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidCreds    = errors.New("invalid credentials")
	ErrRefreshNotFound = errors.New("refresh token not found")
)

type Service struct {
	users   UserRepo
	rtRepo  RefreshTokenRepo
	jwt     *jwtinfra.Service
	nowFunc func() time.Time
}

func NewService(u UserRepo, r RefreshTokenRepo, jwt *jwtinfra.Service) *Service {
	return &Service{
		users:   u,
		rtRepo:  r,
		jwt:     jwt,
		nowFunc: time.Now,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (*TokenPair, error) {
	// existence check
	if u, _ := s.users.GetByEmail(ctx, email); u != nil {
		return nil, ErrUserExists
	}

	now := s.nowFunc().UTC()
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	u := &domainuser.User{
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}

	access, refresh, err := s.jwt.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	rt := &domaintoken.RefreshToken{
		UserID:    u.ID,
		Token:     refresh,
		ExpiresAt: s.nowFunc().Add(s.jwt.RefreshTTL()).Unix(),
	}
	if err := s.rtRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil || u == nil {
		return nil, ErrInvalidCreds
	}
	if !crypto.ComparePassword(u.PasswordHash, password) {
		return nil, ErrInvalidCreds
	}

	access, refresh, err := s.jwt.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	rt := &domaintoken.RefreshToken{
		UserID:    u.ID,
		Token:     refresh,
		ExpiresAt: s.nowFunc().Add(s.jwt.RefreshTTL()).Unix(),
	}
	if err := s.rtRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) Refresh(ctx context.Context, oldRefresh string) (*TokenPair, error) {
	uid, err := s.jwt.ParseRefresh(oldRefresh)
	if err != nil {
		return nil, err
	}

	ok, err := s.rtRepo.Exists(ctx, oldRefresh)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrRefreshNotFound
	}

	// rotate: delete old refresh
	if err := s.rtRepo.Delete(ctx, oldRefresh); err != nil {
		return nil, err
	}

	access, refresh, err := s.jwt.Generate(uid)
	if err != nil {
		return nil, err
	}

	rt := &domaintoken.RefreshToken{
		UserID:    uid,
		Token:     refresh,
		ExpiresAt: s.nowFunc().Add(s.jwt.RefreshTTL()).Unix(),
	}
	if err := s.rtRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) Logout(ctx context.Context, refresh string) error {
	return s.rtRepo.Delete(ctx, refresh)
}

func (s *Service) Me(ctx context.Context, id int64) (*domainuser.User, error) {
	return s.users.GetByID(ctx, id)
}
