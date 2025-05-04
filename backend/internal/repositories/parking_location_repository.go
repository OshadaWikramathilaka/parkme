package repositories

import (
	"context"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParkingLocationRepository struct {
	collection *qmgo.Collection
}

func NewParkingLocationRepository(db *qmgo.Database) *ParkingLocationRepository {
	return &ParkingLocationRepository{
		collection: db.Collection("parking_locations"),
	}
}

func (r *ParkingLocationRepository) Create(ctx context.Context, location *models.ParkingLocation) error {
	_, err := r.collection.InsertOne(ctx, location)
	return err
}

func (r *ParkingLocationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.ParkingLocation, error) {
	var location models.ParkingLocation
	err := r.collection.Find(ctx, bson.M{"_id": id}).One(&location)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &location, nil
}

func (r *ParkingLocationRepository) FindAll(ctx context.Context) ([]models.ParkingLocation, error) {
	var locations []models.ParkingLocation
	err := r.collection.Find(ctx, bson.M{}).All(&locations)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *ParkingLocationRepository) Update(ctx context.Context, location *models.ParkingLocation) error {
	return r.collection.UpdateOne(ctx, bson.M{"_id": location.ID}, bson.M{"$set": location})
}

func (r *ParkingLocationRepository) UpdateSlotStatus(ctx context.Context, locationID primitive.ObjectID, slotNumber string, isOccupied bool) error {
	return r.collection.UpdateOne(ctx,
		bson.M{"_id": locationID, "slots.number": slotNumber},
		bson.M{"$set": bson.M{"slots.$.is_occupied": isOccupied}})
}

func (r *ParkingLocationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	err := r.collection.Remove(ctx, bson.M{"_id": id})
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}
