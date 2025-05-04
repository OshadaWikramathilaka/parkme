package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" `
	PlateNumber string             `json:"plate_number" bson:"plate_number" validate:"required"`
	Brand       string             `json:"brand" bson:"brand" validate:"required"`
	Model       string             `json:"model" bson:"model" validate:"required"`
	Owner       primitive.ObjectID `json:"owner" bson:"owner" validate:"required"`
}

func (v Vehicle) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.PlateNumber, validation.Required),
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Owner, validation.Required),
	)
}
