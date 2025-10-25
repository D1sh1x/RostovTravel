package service

import (
	"hackaton/internal/repository"

	"github.com/rs/zerolog"
)

type service struct {
	user *userService
}

func NewService(repo repository.RepositoryInterface, jwtSecret []byte, logger zerolog.Logger) ServiceInterface {
	return &service{
		user: newUserService(repo, jwtSecret, logger),
	}
}

func (s *service) User() UserServiceInterface { return s.user }
