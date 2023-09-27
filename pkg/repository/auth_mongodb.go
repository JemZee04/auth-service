package repository

import (
	"auth-service/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AuthMongodb struct {
	c *mongo.Collection
}

func NewAuthMongodb(ctx context.Context, client *mongo.Client) (*AuthMongodb, error) {
	dao := &AuthMongodb{
		c: client.Database("auth-service").Collection(usersTable),
	}
	err := dao.createIndices(ctx)
	if err != nil {
		return nil, err
	}
	return dao, nil
}

func (r *AuthMongodb) createIndices(ctx context.Context) error {
	_, err := r.c.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{"expireAt", 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	)
	return err
}

func (r *AuthMongodb) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.c.InsertOne(ctx, user)
	return err
}

func (r *AuthMongodb) GetUser(ctx context.Context, email, password string) (model.User, error) {
	var user model.User
	if err := r.c.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *AuthMongodb) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	var user model.User
	if err := r.c.FindOne(
		ctx, bson.M{
			"session.refreshToken": refreshToken,
			"session.expiresAt":    bson.M{"$gt": time.Now()},
		},
	).Decode(&user); err != nil {

		return model.User{}, err
	}

	return user, nil
}

func (r *AuthMongodb) SetSession(ctx context.Context, userID primitive.ObjectID, session model.Session) error {
	_, err := r.c.UpdateOne(
		ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}},
	)

	return err
}
