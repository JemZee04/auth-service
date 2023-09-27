package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	usersTable = "user"
)

type Config struct {
	Host string
	Port string
}

func NewMongoDB(ctx context.Context, cfg Config) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)))
	return client, err
}
