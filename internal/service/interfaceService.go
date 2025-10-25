package service

import (
	"context"
	"hackaton/internal/dto/request"
	"hackaton/internal/dto/response"
)

type UserServiceInterface interface {
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
	CreateUser(ctx context.Context, req *request.UserRequest) error
	GetUsers(ctx context.Context) ([]*response.UserResponse, error)
	GetUserByID(ctx context.Context, userID string) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, req *request.UserRequest) error
	DeleteUser(ctx context.Context, userID string) error
}

type ServiceInterface interface {
	User() UserServiceInterface
}
