package service

import (
	"auth-service/model"
	"auth-service/pkg/repository"
	"context"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, email, password string) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
