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

type VehicleController struct {
	service *services.VehicleService
}

func NewVehicleController(service *services.VehicleService) *VehicleController {
	return &VehicleController{
		service: service,
	}
}

func (c *VehicleController) GetAll(ctx echo.Context) error {
	vehicles, err := c.service.GetAll(ctx.Request().Context())
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get vehicles", err)
	}
	return utils.SuccessResponse(ctx, http.StatusOK, "Vehicles retrieved successfully", vehicles)
}

func (c *VehicleController) GetByID(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	vehicle, err := c.service.GetByID(ctx.Request().Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Vehicle not found", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Vehicle retrieved successfully", vehicle)
}

func (c *VehicleController) GetUserVehicles(ctx echo.Context) error {
	userID, err := primitive.ObjectIDFromHex(ctx.Param("userId"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID format", err)
	}

	vehicles, err := c.service.GetByOwner(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user vehicles", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "User vehicles retrieved successfully", vehicles)
}

func (c *VehicleController) Create(ctx echo.Context) error {
	var vehicle models.Vehicle
	body, _ := io.ReadAll(ctx.Request().Body)
	fmt.Printf("Raw Request Body: %s\n", string(body))
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	userID := ctx.Get("userID").(primitive.ObjectID)

	vehicle.Owner = userID

	if err := ctx.Bind(&vehicle); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the vehicle
	if err := vehicle.Validate(); err != nil {
		if e, ok := err.(validation.Errors); ok {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", e)
		}
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid vehicle data", err)
	}

	if err := c.service.Create(ctx.Request().Context(), &vehicle); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create vehicle", err)
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "Vehicle created successfully", vehicle)
}

func (c *VehicleController) Update(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	userID := ctx.Get("userID").(primitive.ObjectID)

	var vehicle models.Vehicle
	if err := ctx.Bind(&vehicle); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}
	vehicle.Owner = userID
	vehicle.ID = id

	// Validate the vehicle
	if err := vehicle.Validate(); err != nil {
		if e, ok := err.(validation.Errors); ok {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", e)
		}
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid vehicle data", err)
	}

	if err := c.service.Update(ctx.Request().Context(), &vehicle); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update vehicle", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Vehicle updated successfully", vehicle)
}

func (c *VehicleController) Delete(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err)
	}

	if err := c.service.Delete(ctx.Request().Context(), id); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete vehicle", err)
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Vehicle deleted successfully", nil)
}
