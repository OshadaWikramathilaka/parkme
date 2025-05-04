package controllers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) GetUserService() *services.UserService {
	return c.service
}

func (c *UserController) GetAll(ctx echo.Context) error {
	users, err := c.service.GetAll(ctx.Request().Context())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get users", err)
	}
	for i := range users {
		users[i].Password = "" // Do not return password in response
	}
	return utils.SuccessResponse(ctx, http.StatusOK, "Users retrieved successfully", users)
}

func (c *UserController) GetByID(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	user, err := c.service.GetByID(ctx.Request().Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", err)
	}
	user.Password = "" // Do not return password in response

	return utils.SuccessResponse(ctx, http.StatusOK, "User retrieved successfully", user)
}

func (c *UserController) Create(ctx echo.Context) error {
	var user models.User
	body, _ := io.ReadAll(ctx.Request().Body)
	fmt.Printf("Raw Request Body: %s\n", string(body))
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := ctx.Bind(&user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	// Call BeforeCreate which includes validation
	if err := user.BeforeCreate(); err != nil {
		// Handle validation errors specifically
		if e, ok := err.(validation.Errors); ok {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", e)
		}
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user data", err)
	}

	if err := c.service.Create(ctx.Request().Context(), &user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user", err)
	}
	user.Password = "" // Do not return password in response

	return utils.SuccessResponse(ctx, http.StatusCreated, "User created successfully", user)
}

func (c *UserController) Update(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	user.ID = id

	// Call BeforeUpdate which includes validation
	if err := user.BeforeUpdate(); err != nil {
		// Handle validation errors specifically
		if e, ok := err.(validation.Errors); ok {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", e)
		}
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user data", err)
	}

	if err := c.service.Update(ctx.Request().Context(), &user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update user", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "User updated successfully", user)
}

func (c *UserController) Delete(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete user", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "User deleted successfully", nil)
}
