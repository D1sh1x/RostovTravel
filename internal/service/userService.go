package service

import (
	"context"
	"errors"
	"fmt"

	"hackaton/internal/dto/request"
	"hackaton/internal/dto/response"
	"hackaton/internal/models"
	"hackaton/internal/repository"
	"hackaton/internal/utils/jwt"
	rolevalidate "hackaton/internal/utils/roleValidate"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo      repository.RepositoryInterface
	jwtSecret []byte
	logger    zerolog.Logger
}

func newUserService(repo repository.RepositoryInterface, jwtSecret []byte, logger zerolog.Logger) *userService {
	return &userService{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	log := s.logger.With().Str("name", req.Name).Logger()
	log.Info().Msg("login attempt")

	user, err := s.repo.User().GetUserByName(ctx, req.Name)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by name")
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn().Msg("invalid password")
		return nil, errors.New("invalid credentials")
	}

	role := user.Role
	if role == "" {
		role = "user"
	}

	token, err := jwt.GenerateToken(user.ID.Hex(), user.Name, role, s.jwtSecret)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Info().Msg("login successful")

	return &response.LoginResponse{
		Token: token,
		Role:  role,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *request.UserRequest) error {
	log := s.logger.With().Str("name", req.Name).Logger()
	log.Info().Msg("creating user")

	if existing, err := s.repo.User().GetUserByName(ctx, req.Name); err == nil && existing != nil {
		log.Warn().Msg("user already exists")
		return errors.New("user already exists")
	}

	if req.Password == "" {
		log.Warn().Msg("empty password")
		return errors.New("password is empty")
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return err
	}

	favoritesModel := make([]models.Favorite, 0, len(req.Favorites))
	for _, f := range req.Favorites {
		favoritesModel = append(favoritesModel, models.Favorite{
			Type:   f.Type,
			ItemID: f.ItemID,
		})
	}

	u := &models.User{
		Name:         req.Name,
		PasswordHash: string(passwordHashBytes),
		Favorites:    favoritesModel,
		Role:         req.Role,
	}

	if err := s.repo.User().CreateUser(ctx, u); err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}

	log.Info().Str("id", u.ID.Hex()).Msg("user created")
	return nil
}

func (s *userService) GetUsers(ctx context.Context) ([]*response.UserResponse, error) {
	log := s.logger.With().Logger()
	log.Info().Msg("getting all users")

	users, err := s.repo.User().GetUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get users")
		return nil, err
	}

	result := make([]*response.UserResponse, 0, len(users))
	for _, u := range users {
		favs := make([]response.FavoriteResponse, 0, len(u.Favorites))
		for _, f := range u.Favorites {
			favs = append(favs, response.FavoriteResponse{
				Type:   f.Type,
				ItemID: f.ItemID,
			})
		}

		result = append(result, &response.UserResponse{
			ID:           u.ID.Hex(),
			Name:         u.Name,
			PasswordHash: u.PasswordHash,
			CreatedAt:    u.CreatedAt,
			Favorites:    favs,
			Role:         u.Role,
		})
	}

	return result, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*response.UserResponse, error) {
	log := s.logger.With().Str("id", userID).Logger()
	log.Info().Msg("getting user")

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Warn().Err(err).Msg("invalid user id format")
		return nil, fmt.Errorf("invalid user id")
	}

	u, err := s.repo.User().GetUserByID(ctx, oid)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by id")
		return nil, err
	}

	favs := make([]response.FavoriteResponse, 0, len(u.Favorites))
	for _, f := range u.Favorites {
		favs = append(favs, response.FavoriteResponse{
			Type:   f.Type,
			ItemID: f.ItemID,
		})
	}

	resp := &response.UserResponse{
		ID:           u.ID.Hex(),
		Name:         u.Name,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		Favorites:    favs,
		Role:         u.Role,
	}

	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, userID string, req *request.UserRequest) error {
	log := s.logger.With().Str("id", userID).Logger()
	log.Info().Msg("updating user")

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Warn().Err(err).Msg("invalid user id format")
		return fmt.Errorf("invalid user id")
	}

	u, err := s.repo.User().GetUserByID(ctx, oid)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user before update")
		return err
	}

	if req.Name != "" {
		u.Name = req.Name
	}

	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error().Err(err).Msg("failed to hash password on update")
			return err
		}
		u.PasswordHash = string(hashed)
	}

	if req.Favorites != nil {
		newFavs := make([]models.Favorite, 0, len(req.Favorites))
		for _, f := range req.Favorites {
			newFavs = append(newFavs, models.Favorite{
				Type:   f.Type,
				ItemID: f.ItemID,
			})
		}
		u.Favorites = newFavs
	}

	if req.Role != "" {
		if !rolevalidate.IsValidRole(req.Role) {
			log.Warn().Str("role", req.Role).Msg("invalid role")
			return errors.New("invalid role")
		}
		u.Role = req.Role
	}

	if err := s.repo.User().UpdateUser(ctx, u); err != nil {
		log.Error().Err(err).Msg("failed to update user")
		return err
	}

	log.Info().Str("id", userID).Msg("user updated")
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	log := s.logger.With().Str("id", userID).Logger()
	log.Info().Msg("deleting user")

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Warn().Err(err).Msg("invalid user id format")
		return fmt.Errorf("invalid user id")
	}

	if err := s.repo.User().DeleteUser(ctx, oid); err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		return err
	}

	log.Info().Msg("user deleted")
	return nil
}
