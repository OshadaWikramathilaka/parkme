package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"

	MinPasswordLength = 8
	MaxPasswordLength = 72
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password,omitempty" bson:"password"`
	Role            Role               `json:"role" bson:"role"`
	Status          string             `json:"status" bson:"status"`
	ProfileImageURL string             `json:"profile_image_url,omitempty" bson:"profile_image_url,omitempty"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		// Name validation
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 50).Error("name must be between 2 and 50 characters"),
		),

		// Email validation
		validation.Field(&u.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("invalid email format"),
		),

		// Password validation
		validation.Field(&u.Password,
			validation.Required.Error("password is required"),
			validation.Length(MinPasswordLength, MaxPasswordLength).
				Error("password must be between 8 and 72 characters"),
		),

		// Role validation
		validation.Field(&u.Role,
			validation.Required.Error("role is required"),
			validation.In(RoleAdmin, RoleUser).Error("invalid role"),
		),

		// Status validation
		validation.Field(&u.Status,
			validation.Required.Error("status is required"),
			validation.In(UserStatusActive, UserStatusInactive, UserStatusBanned).
				Error("invalid status"),
		),
	)
}

func (u *User) ValidateUpdate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.Required.Error("name is required"),
			validation.Length(2, 50).Error("name must be between 2 and 50 characters"),
		),
		validation.Field(&u.Email,
			validation.Required.Error("email is required"),
			is.Email.Error("invalid email format"),
		),
		// Password is optional during update
		validation.Field(&u.Password,
			validation.When(len(u.Password) > 0, validation.Length(MinPasswordLength, MaxPasswordLength).
				Error("password must be between 8 and 72 characters")),
		),
		validation.Field(&u.Role,
			validation.Required.Error("role is required"),
			validation.In(RoleAdmin, RoleUser).Error("invalid role"),
		),
		validation.Field(&u.Status,
			validation.Required.Error("status is required"),
			validation.In(UserStatusActive, UserStatusInactive, UserStatusBanned).
				Error("invalid status"),
		),
	)
}

func (u *User) HashPassword() error {
	if len(u.Password) == 0 {
		return validation.NewError("validation_error", "password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	// Add null check for hashed password
	if u.Password == "" {
		return errors.New("password not set")
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func (u *User) BeforeCreate() error {
	// Set timestamps
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	// Set defaults if empty
	if u.Role == "" {
		u.Role = RoleUser
	}
	if u.Status == "" {
		u.Status = UserStatusActive
	}

	// Validate all fields
	if err := u.Validate(); err != nil {
		return err
	}

	// Hash password
	return u.HashPassword()
}

func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()

	// Validate fields
	if err := u.ValidateUpdate(); err != nil {
		return err
	}

	// Hash password if provided
	if u.Password != "" {
		return u.HashPassword()
	}

	return nil
}
