package repository

import (
	"auth-service/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, email, password string) (model.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session model.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(ctx context.Context, db *mongo.Client) *Repository {
	auth, err := NewAuthMongodb(ctx, db)
	if err != nil {
		return nil
	}
	return &Repository{
		Authorization: auth,
	}
}
