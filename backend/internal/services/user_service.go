package services

import (
	"context"
	"errors"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) FindOne(ctx context.Context, filter interface{}) (*models.User, error) {
	return s.repo.FindOne(ctx, filter)
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	// Check if the user already exists
	filter := bson.M{"email": user.Email}
	existingUser, err := s.FindOne(ctx, filter)
	if err != nil && err != repositories.ErrNotFound {
		return err
	}
	if existingUser != nil {
		return ErrUserExists
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
