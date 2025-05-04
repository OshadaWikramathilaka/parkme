package controllers

import (
	"net/http"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingController struct {
	service *services.BookingService
}

func NewBookingController(service *services.BookingService) *BookingController {
	return &BookingController{service: service}
}

func (c *BookingController) CreateBooking(ctx echo.Context) error {
	var booking models.Booking
	if err := ctx.Bind(&booking); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userID := ctx.Get("userID").(primitive.ObjectID)
	booking.UserID = userID

	if err := c.service.CreateBooking(ctx.Request().Context(), &booking); err != nil {
		switch err {
		case services.ErrInvalidTimeRange:
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid time range")
		case services.ErrSpotAlreadyBooked:
			return echo.NewHTTPError(http.StatusConflict, "Spot already booked")
		case services.ErrBookingInPast:
			return echo.NewHTTPError(http.StatusBadRequest, "Cannot book in the past")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "Booking created successfully", booking)
}

func (c *BookingController) GetBooking(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid booking ID")
	}

	booking, err := c.service.GetBooking(ctx.Request().Context(), id)
	if err != nil {
		if err == services.ErrBookingNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Booking not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Booking retrieved successfully", booking)
}

func (c *BookingController) GetVehicleBookings(ctx echo.Context) error {
	vehicleID, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid vehicle ID")
	}

	bookings, err := c.service.GetVehicleBookings(ctx.Request().Context(), vehicleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Vehicle bookings retrieved successfully", bookings)
}

func (c *BookingController) CancelBooking(ctx echo.Context) error {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid booking ID")
	}

	if err := c.service.CancelBooking(ctx.Request().Context(), id); err != nil {
		if err == services.ErrBookingNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Booking not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

func (c *BookingController) GetUserBookings(ctx echo.Context) error {
	userID := ctx.Get("userID").(primitive.ObjectID)
	bookings, err := c.service.GetUserBookings(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user bookings", err)
	}
	return utils.SuccessResponse(ctx, http.StatusOK, "User bookings retrieved successfully", bookings)
}
