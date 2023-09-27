package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" binding:"required" bson:"email"`
	Phone    string             `json:"phone" binding:"required" bson:"phone"`
	Password string             `json:"password" binding:"required" bson:"password"`
}
