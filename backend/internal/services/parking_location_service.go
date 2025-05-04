package services

import (
	"context"
	"errors"
	"time"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrLocationNotFound = errors.New("parking location not found")
	ErrSlotNotFound     = errors.New("parking slot not found")
)

type ParkingLocationService struct {
	repo *repositories.ParkingLocationRepository
}

func NewParkingLocationService(repo *repositories.ParkingLocationRepository) *ParkingLocationService {
	return &ParkingLocationService{
		repo: repo,
	}
}

func (s *ParkingLocationService) CreateLocation(ctx context.Context, location *models.ParkingLocation) error {
	location.CreatedAt = time.Now()
	location.UpdatedAt = time.Now()
	location.TotalSlots = len(location.Slots)
	return s.repo.Create(ctx, location)
}

func (s *ParkingLocationService) GetLocation(ctx context.Context, id primitive.ObjectID) (*models.ParkingLocation, error) {
	location, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	return location, nil
}

func (s *ParkingLocationService) GetAllLocations(ctx context.Context) ([]models.ParkingLocation, error) {
	return s.repo.FindAll(ctx)
}

func (s *ParkingLocationService) UpdateLocation(ctx context.Context, location *models.ParkingLocation) error {
	location.UpdatedAt = time.Now()
	location.TotalSlots = len(location.Slots)
	return s.repo.Update(ctx, location)
}

func (s *ParkingLocationService) UpdateSlotStatus(ctx context.Context, locationID primitive.ObjectID, slotNumber string, isOccupied bool) error {
	// First verify the location and slot exist
	location, err := s.GetLocation(ctx, locationID)
	if err != nil {
		return err
	}

	slotExists := false
	for _, slot := range location.Slots {
		if slot.Number == slotNumber {
			slotExists = true
			break
		}
	}

	if !slotExists {
		return ErrSlotNotFound
	}

	return s.repo.UpdateSlotStatus(ctx, locationID, slotNumber, isOccupied)
}

func (s *ParkingLocationService) DeleteLocation(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
