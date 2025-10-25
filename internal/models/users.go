package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	PasswordHash string             `bson:"password_hash" json:"password_hash"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	Favorites    []Favorite         `bson:"favorites" json:"favorites"`
	Role         string             `bson:"role" json:"role"`
	Username     string             `bson:"username" json:"username"`
}

type Favorite struct {
	Type   string `bson:"type" json:"type"`
	ItemID string `bson:"item_id" json:"item_id"`
}
