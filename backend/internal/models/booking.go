package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusActive    BookingStatus = "active"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

type BookingType string

const (
	BookingTypePreBooked BookingType = "pre_booked"
	BookingTypeOnSite    BookingType = "on_site"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	VehicleID   primitive.ObjectID `bson:"vehicle_id" json:"vehicleId"`
	UserID      primitive.ObjectID `bson:"user_id" json:"userId"`
	LocationID  primitive.ObjectID `bson:"location_id" json:"locationId"`
	StartTime   time.Time          `bson:"start_time" json:"startTime"`
	EndTime     *time.Time         `bson:"end_time,omitempty" json:"endTime,omitempty"`
	Status      BookingStatus      `bson:"status" json:"status"`
	SpotNumber  *string            `bson:"spot_number,omitempty" json:"spotNumber,omitempty"`
	TotalAmount *float64           `bson:"total_amount,omitempty" json:"totalAmount,omitempty"`
	BookingType BookingType        `bson:"booking_type" json:"bookingType"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	User        *User              `bson:"-" json:"user,omitempty"`
	Vehicle     *Vehicle           `bson:"-" json:"vehicle,omitempty"`
	Location    *ParkingLocation   `bson:"-" json:"location,omitempty"`
}
