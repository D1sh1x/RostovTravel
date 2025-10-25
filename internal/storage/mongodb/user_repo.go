package mongodb

import (
	"context"
	"fmt"
	"hackaton/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepo struct {
	coll *mongo.Collection
}

func (r *UserMongoRepo) CreateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now

	res, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Code == 11000 {
					return fmt.Errorf("duplicate username")
				}
			}
		}
		if ce, ok := err.(mongo.CommandError); ok {
			if ce.Code == 11000 {
				return fmt.Errorf("duplicate username")
			}
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserMongoRepo) GetUsers(ctx context.Context) ([]*models.User, error) {
	cur, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cur.Close(ctx)

	var users []*models.User
	if err = cur.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

func (r *UserMongoRepo) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserMongoRepo) GetUserByName(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}
	return &user, nil
}

func (r *UserMongoRepo) UpdateUser(ctx context.Context, user *models.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":          user.Name,
			"password_hash": user.PasswordHash,
			"favorites":     user.Favorites,
			"role":          user.Role,
		},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *UserMongoRepo) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
