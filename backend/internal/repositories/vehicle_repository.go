package repositories

import (
	"context"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleRepository struct {
	collection *qmgo.Collection
}

func NewVehicleRepository(db *qmgo.Database) *VehicleRepository {
	return &VehicleRepository{
		collection: db.Collection("vehicles"),
	}
}

func (r *VehicleRepository) FindOne(ctx context.Context, filter interface{}) (*models.Vehicle, error) {
	var vehicle models.Vehicle

	err := r.collection.Find(ctx, filter).One(&vehicle)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) FindAll(ctx context.Context) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.collection.Find(ctx, bson.M{}).All(&vehicles)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	err := r.collection.Find(ctx, bson.M{"_id": id}).One(&vehicle)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &vehicle, nil
}

func (r *VehicleRepository) FindByOwner(ctx context.Context, ownerID primitive.ObjectID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.collection.Find(ctx, bson.M{"owner": ownerID}).All(&vehicles)
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (r *VehicleRepository) Create(ctx context.Context, vehicle *models.Vehicle) error {
	_, err := r.collection.InsertOne(ctx, vehicle)

	return err
}

func (r *VehicleRepository) Update(ctx context.Context, vehicle *models.Vehicle) error {
	err := r.collection.UpdateOne(ctx, bson.M{"_id": vehicle.ID}, bson.M{"$set": vehicle})
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *VehicleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	err := r.collection.RemoveId(ctx, id)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}
