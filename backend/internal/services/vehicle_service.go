package services

import (
	"context"
	"errors"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrVehicleExists   = errors.New("vehicle with this plate number already exists")
	ErrVehicleNotFound = errors.New("vehicle not found")
	ErrInvalidOwner    = errors.New("invalid vehicle owner")
)

type VehicleService struct {
	repo *repositories.VehicleRepository
}

func NewVehicleService(repo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{
		repo: repo,
	}
}

func (s *VehicleService) GetAll(ctx context.Context) ([]models.Vehicle, error) {
	return s.repo.FindAll(ctx)
}

func (s *VehicleService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Vehicle, error) {
	vehicle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrVehicleNotFound
		}
		return nil, err
	}
	return vehicle, nil
}

func (s *VehicleService) GetByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Vehicle, error) {
	return s.repo.FindByOwner(ctx, ownerID)
}

func (s *VehicleService) Create(ctx context.Context, vehicle *models.Vehicle) error {
	// Check if vehicle with same plate number exists
	filter := bson.M{"plate_number": vehicle.PlateNumber}
	existing, err := s.repo.FindOne(ctx, filter)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return err
	}
	if existing != nil {
		return ErrVehicleExists
	}

	return s.repo.Create(ctx, vehicle)
}

func (s *VehicleService) Update(ctx context.Context, vehicle *models.Vehicle) error {
	// Check if vehicle exists
	existing, err := s.repo.FindByID(ctx, vehicle.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrVehicleNotFound
		}
		return err
	}

	// Check if updating to a plate number that already exists on another vehicle
	if vehicle.PlateNumber != existing.PlateNumber {
		filter := bson.M{
			"plate_number": vehicle.PlateNumber,
			"_id":          bson.M{"$ne": vehicle.ID},
		}
		duplicate, err := s.repo.FindOne(ctx, filter)
		if err != nil && !errors.Is(err, repositories.ErrNotFound) {
			return err
		}
		if duplicate != nil {
			return ErrVehicleExists
		}
	}

	return s.repo.Update(ctx, vehicle)
}

func (s *VehicleService) Delete(ctx context.Context, id primitive.ObjectID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return ErrVehicleNotFound
		}
		return err
	}
	return nil
}

func (s *VehicleService) FindOne(ctx context.Context, filter bson.M) (*models.Vehicle, error) {
	return s.repo.FindOne(ctx, filter)
}
