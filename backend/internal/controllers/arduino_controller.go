package controllers

import (
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArduinoController struct {
	service         *services.ArduinoService
	bookingService  *services.BookingService
	userService     *services.UserService
	locationService *services.ParkingLocationService
	walletService   *services.WalletService
}

func NewArduinoController(
	service *services.ArduinoService,
	bookingService *services.BookingService,
	userService *services.UserService,
	locationService *services.ParkingLocationService,
	walletService *services.WalletService,
) *ArduinoController {
	return &ArduinoController{
		service:         service,
		bookingService:  bookingService,
		userService:     userService,
		locationService: locationService,
		walletService:   walletService,
	}
}

type GateEnterRequest struct {
	LocationID string                `json:"location_id" form:"location_id"`
	Image      *multipart.FileHeader `json:"image" form:"image"`
}

// get the upload image and extract the vehicle number plate from it
func (c *ArduinoController) GateEnter(ctx echo.Context) error {
	var req GateEnterRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err)
	}

	// Validate location ID
	if req.LocationID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Location ID is required", nil)
	}

	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid location ID format", err)
	}

	// Get the image from the request
	if req.Image == nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Image file is required", nil)
	}
	file, err := req.Image.Open()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to read image file", err)
	}
	defer file.Close()

	ext := filepath.Ext(req.Image.Filename)

	imageName := uuid.New().String() + ext

	// Save image to folder first
	_, err = c.service.SaveImageToFolder(file, imageName)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save image to folder", err)
	}

	// Reopen the file for the plate extraction
	file, err = req.Image.Open()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to reopen image file", err)
	}
	defer file.Close()

	// Extract image data and validate vehicle
	plateResult, err := c.service.ExtractNumberPlate(file)
	if err != nil {
		if err == services.ErrVehicleNotFound {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Vehicle not registered in the system", err)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to extract image data", err)
	}

	// Check if vehicle already has an active booking
	existingBooking, err := c.bookingService.FindBookingByFilter(ctx.Request().Context(), bson.M{
		"vehicle_id": plateResult.Vehicle.ID,
		"status": bson.M{
			"$in": []string{
				string(models.BookingStatusActive),
				string(models.BookingStatusPending),
			},
		},
	})
	if err != nil && err != services.ErrBookingNotFound {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to check existing bookings", err)
	}
	if existingBooking != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, "Vehicle already has an active booking", fmt.Errorf("duplicate booking not allowed"))
	}

	// Get owner details
	owner, err := c.userService.GetByID(ctx.Request().Context(), plateResult.Vehicle.Owner)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Vehicle owner not found", err)
	}

	// Check for existing pre-booked booking
	existingBooking, err = c.bookingService.FindBookingByFilter(ctx.Request().Context(), bson.M{
		"vehicle_id":   plateResult.Vehicle.ID,
		"location_id":  locationID,
		"status":       models.BookingStatusPending,
		"booking_type": models.BookingTypePreBooked,
		"start_time": bson.M{
			"$lte": time.Now(),
		},
		"end_time": bson.M{
			"$gte": time.Now(),
		},
	})

	var booking *models.Booking
	if err == nil && existingBooking != nil {
		// Update existing pre-booked booking
		existingBooking.Status = models.BookingStatusActive
		if err := c.bookingService.UpdateBookingStatus(ctx.Request().Context(), existingBooking.ID, models.BookingStatusActive); err != nil {
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update existing booking", err)
		}
		booking = existingBooking
	} else {
		// Create new on-site booking
		now := time.Now()
		booking = &models.Booking{
			VehicleID:   plateResult.Vehicle.ID,
			UserID:      plateResult.Vehicle.Owner,
			LocationID:  locationID,
			StartTime:   now,
			Status:      models.BookingStatusActive, // Set as active immediately
			BookingType: models.BookingTypeOnSite,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// For on-site bookings at the gate, we don't specify a spot number
		// Let the booking service find an available spot
		if err := c.bookingService.CreateBooking(ctx.Request().Context(), booking); err != nil {
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create booking", err)
		}
	}

	// Clear sensitive owner information
	owner.Password = ""

	// Return success response with plate result, vehicle, owner, and booking
	// response := map[string]interface{}{
	// 	"plateResult": plateResult,
	// 	"vehicle":     plateResult.Vehicle,
	// 	"owner":       owner,
	// 	"booking":     booking,
	// }
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Vehicle processed successfully",
	})
}

type SpotData struct {
	LocationID string `json:"location_id"`
	SpotNumber string `json:"spot_number"`
	IsOccupied bool   `json:"is_occupied"`
}

// UpdateSpotStatus updates the spot status based on sensor data from Arduino
func (c *ArduinoController) UpdateSpotStatus(ctx echo.Context) error {
	var spotData SpotData
	if err := ctx.Bind(&spotData); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request data", err)
	}

	// Validate location ID
	if spotData.LocationID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Location ID is required", nil)
	}

	locationID, err := primitive.ObjectIDFromHex(spotData.LocationID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid location ID format", err)
	}

	// Validate spot number
	if spotData.SpotNumber == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Spot number is required", nil)
	}

	// Update spot status in the location
	err = c.locationService.UpdateSlotStatus(ctx.Request().Context(), locationID, spotData.SpotNumber, spotData.IsOccupied)
	if err != nil {
		switch err {
		case services.ErrLocationNotFound:
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Location not found", err)
		case services.ErrSlotNotFound:
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Spot not found", err)
		default:
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update spot status", err)
		}
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Spot status updated successfully", map[string]interface{}{
		"location_id": locationID,
		"spot_number": spotData.SpotNumber,
		"is_occupied": spotData.IsOccupied,
	})
}

type GateExitRequest struct {
	LocationID string                `json:"location_id" form:"location_id"`
	Image      *multipart.FileHeader `json:"image" form:"image"`
}

// GateExit handles vehicle exit, calculates payment, and updates spot status
func (c *ArduinoController) GateExit(ctx echo.Context) error {
	var req GateExitRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request", err)
	}

	// Validate location ID
	if req.LocationID == "" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Location ID is required", nil)
	}

	locationID, err := primitive.ObjectIDFromHex(req.LocationID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid location ID format", err)
	}

	// Get the image from the request
	if req.Image == nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Image file is required", nil)
	}
	file, err := req.Image.Open()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to read image file", err)
	}
	defer file.Close()

	ext := filepath.Ext(req.Image.Filename)

	imageName := uuid.New().String() + "_exit" + ext

	// Save image to folder first
	_, err = c.service.SaveImageToFolder(file, imageName)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save image to folder", err)
	}

	// Reopen the file for the plate extraction
	file, err = req.Image.Open()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to reopen image file", err)
	}
	defer file.Close()

	// Extract image data and validate vehicle
	plateResult, err := c.service.ExtractNumberPlate(file)
	if err != nil {
		if err == services.ErrVehicleNotFound {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "Vehicle not registered in the system", err)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to extract image data", err)
	}

	// Find active booking for this vehicle
	activeBooking, err := c.bookingService.FindBookingByFilter(ctx.Request().Context(), bson.M{
		"vehicle_id": plateResult.Vehicle.ID,
		"status": bson.M{
			"$in": []string{
				string(models.BookingStatusActive),
				string(models.BookingStatusPending),
			},
		},
	})
	if err != nil {
		if err == services.ErrBookingNotFound {
			return utils.ErrorResponse(ctx, http.StatusNotFound, "No active booking found for this vehicle", err)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find active booking", err)
	}

	// Calculate parking duration and amount
	endTime := time.Now()
	duration := endTime.Sub(activeBooking.StartTime)
	hours := math.Ceil(duration.Hours())
	totalAmount := hours * 100 // 100 points per hour

	// Get owner details and wallet
	owner, err := c.userService.GetByID(ctx.Request().Context(), plateResult.Vehicle.Owner)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "Vehicle owner not found", err)
	}

	// Deduct payment from wallet
	transaction, err := c.walletService.Deduct(ctx.Request().Context(), owner.ID, totalAmount, fmt.Sprintf("Parking payment for %s", plateResult.Vehicle.PlateNumber))
	if err != nil {
		if err == services.ErrInsufficientBalance {
			return utils.ErrorResponse(ctx, http.StatusPaymentRequired, "Insufficient wallet balance", err)
		}
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to process payment", err)
	}

	// Complete the booking with payment details
	bookingErr := c.bookingService.CompleteOnSiteBooking(ctx.Request().Context(), activeBooking.ID, endTime, totalAmount)
	if bookingErr != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to complete booking", err)
	}

	// Update spot status to unoccupied
	if activeBooking.SpotNumber != nil {
		err = c.locationService.UpdateSlotStatus(ctx.Request().Context(), locationID, *activeBooking.SpotNumber, false)
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update spot status", err)
		}
	}

	// Save exit image
	_, err = c.service.SaveImageToFolder(file, req.Image.Filename)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save image to folder", err)
	}

	// Clear sensitive owner information
	owner.Password = ""

	// Prepare response with all details
	// response := map[string]interface{}{
	// 	"plate_result": plateResult,
	// 	"booking": map[string]interface{}{
	// 		"id":          activeBooking.ID,
	// 		"startTime":   activeBooking.StartTime,
	// 		"endTime":     endTime,
	// 		"totalAmount": totalAmount,
	// 	},
	// 	"transaction": transaction,
	// }
	// return utils.SuccessResponse(ctx, http.StatusOK, "Vehicle exit processed successfully", response)

	fmt.Print(transaction)

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Vehicle processed successfully",
	})
}
