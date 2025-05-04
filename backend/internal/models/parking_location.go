package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParkingSlot represents a single parking slot in a location
type ParkingSlot struct {
	Number     string `bson:"number" json:"number"`           // Slot number/identifier
	IsOccupied bool   `bson:"is_occupied" json:"is_occupied"` // Current occupancy status
	Type       string `bson:"type" json:"type"`               // e.g., "standard", "handicap", "electric"
}

// ParkingLocation represents a parking facility
type ParkingLocation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`               // Name of the parking location
	Address    string             `bson:"address" json:"address"`         // Physical address
	Slots      []ParkingSlot      `bson:"slots" json:"slots"`             // List of parking slots
	TotalSlots int                `bson:"total_slots" json:"total_slots"` // Total number of slots
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
