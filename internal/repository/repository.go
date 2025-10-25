package repository

import (
	"context"
	"hackaton/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
	GetUserByName(ctx context.Context, name string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type RepositoryInterface interface {
	User() UserRepositoryInterface
}
