package middleware

import (
	"net/http"
	"strings"

	"github.com/dfanso/parkme-backend/config"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func ValidateAPIKey() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cfg := config.Load()
			// Get the API key from the request headers
			apiKey := c.Request().Header.Get("X-API-KEY")

			// Check if the API key is empty
			if apiKey == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "API key is required", nil)
			}

			// Check if the API key is valid
			if !strings.Contains(cfg.ARDUNIO_API_KEY, apiKey) {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid API key", nil)
			}

			return next(c)
		}
	}
}
