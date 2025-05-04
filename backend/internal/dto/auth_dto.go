package dto

import "github.com/dfanso/parkme-backend/internal/models"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}
