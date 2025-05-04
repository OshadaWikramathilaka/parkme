package controllers

import (
	"net/http"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParkingLocationController struct {
	service *services.ParkingLocationService
}

func NewParkingLocationController(service *services.ParkingLocationService) *ParkingLocationController {
	return &ParkingLocationController{
		service: service,
	}
}

func (c *ParkingLocationController) CreateLocation(ctx echo.Context) error {
	userRole := models.Role(ctx.Get("userRole").(string))
	if userRole != models.RoleAdmin {
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": "You are not authorized to create a parking location"})
	}

	var location models.ParkingLocation
	if err := ctx.Bind(&location); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.service.CreateLocation(ctx.Request().Context(), &location); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, location)
}

func (c *ParkingLocationController) GetLocation(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	location, err := c.service.GetLocation(ctx.Request().Context(), id)
	if err != nil {
		if err == services.ErrLocationNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, location)
}

func (c *ParkingLocationController) GetAllLocations(ctx echo.Context) error {
	locations, err := c.service.GetAllLocations(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, locations)
}

func (c *ParkingLocationController) UpdateLocation(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	var location models.ParkingLocation
	if err := ctx.Bind(&location); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	location.ID = id
	if err := c.service.UpdateLocation(ctx.Request().Context(), &location); err != nil {
		if err == services.ErrLocationNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, location)
}

func (c *ParkingLocationController) UpdateSlotStatus(ctx echo.Context) error {
	locationID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid location ID format"})
	}

	type UpdateSlotRequest struct {
		SlotNumber string `json:"slot_number"`
		IsOccupied bool   `json:"is_occupied"`
	}

	var req UpdateSlotRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.service.UpdateSlotStatus(ctx.Request().Context(), locationID, req.SlotNumber, req.IsOccupied); err != nil {
		if err == services.ErrLocationNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}
		if err == services.ErrSlotNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Slot not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Slot status updated successfully"})
}

func (c *ParkingLocationController) DeleteLocation(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	if err := c.service.DeleteLocation(ctx.Request().Context(), id); err != nil {
		if err == services.ErrLocationNotFound {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Location deleted successfully"})
}
