package controllers

import (
	"net/http"

	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsController struct {
	service *services.UserStatsService
}

func NewUserStatsController(service *services.UserStatsService) *UserStatsController {
	return &UserStatsController{
		service: service,
	}
}

// GetUserStats godoc
// @Summary Get user statistics for mobile app home screen
// @Description Returns various statistics about user's bookings and usage
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.UserStatsResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/user/stats [get]
func (c *UserStatsController) GetUserStats(ctx echo.Context) error {
	userID := ctx.Get("userID").(primitive.ObjectID)

	stats, err := c.service.GetUserStats(ctx.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "User statistics retrieved successfully", stats)
}
