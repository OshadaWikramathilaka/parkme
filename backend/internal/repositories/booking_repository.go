package repositories

import (
	"context"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingRepository struct {
	db *qmgo.Database
}

func NewBookingRepository(db *qmgo.Database) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) collection() *qmgo.Collection {
	return r.db.Collection("bookings")
}

func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	_, err := r.collection().InsertOne(ctx, booking)
	return err
}

func (r *BookingRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Booking, error) {
	var booking models.Booking
	err := r.collection().Find(ctx, bson.M{"_id": id}).One(&booking)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByVehicle(ctx context.Context, vehicleID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.collection().Find(ctx, bson.M{"vehicle_id": vehicleID}).All(&bookings)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) FindActive(ctx context.Context) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.collection().Find(ctx, bson.M{"status": models.BookingStatusActive}).All(&bookings)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *BookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	return r.collection().UpdateOne(ctx, bson.M{"_id": booking.ID}, bson.M{"$set": booking})
}

func (r *BookingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	err := r.collection().Remove(ctx, bson.M{"_id": id})
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *BookingRepository) FindOne(ctx context.Context, filter bson.M) (*models.Booking, error) {
	var booking models.Booking
	err := r.collection().Find(ctx, filter).One(&booking)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &booking, nil
}

// Count returns the number of documents matching the filter
func (r *BookingRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection().Find(ctx, filter).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *BookingRepository) FindByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.collection().Find(ctx, bson.M{"user_id": userID}).Sort("-created_at").All(&bookings)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}
