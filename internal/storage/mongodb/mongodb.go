package mongodb

import (
	"context"
	"fmt"
	"time"

	"hackaton/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client   *mongo.Client
	db       *mongo.Database
	userRepo *UserMongoRepo
}

func NewStorage(connectionString, databaseName, collectionName string) (repository.RepositoryInterface, error) {
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(databaseName)
	usersColl := db.Collection(collectionName)

	idxModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	}

	ctxIdx, cancelIdx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelIdx()

	if _, err := usersColl.Indexes().CreateOne(ctxIdx, idxModel); err != nil {
		return nil, fmt.Errorf("failed to create username index: %w", err)
	}

	storage := &Storage{
		client:   client,
		db:       db,
		userRepo: &UserMongoRepo{coll: usersColl},
	}

	return storage, nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *Storage) User() repository.UserRepositoryInterface {
	return s.userRepo
}
