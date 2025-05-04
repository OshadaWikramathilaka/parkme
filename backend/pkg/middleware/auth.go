package middleware

import (
	"net/http"
	"strings"

	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/auth"
	"github.com/dfanso/parkme-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthMiddleware(jwtManager *auth.JWTManager, userService *services.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is missing", nil)
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			}

			claims, err := jwtManager.ValidateToken(bearerToken[1])
			if err != nil {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", err)
			}

			//check user is valid
			filter := bson.M{"_id": claims.UserID}
			_, err = userService.FindOne(c.Request().Context(), filter)
			if err != nil {
				return utils.ErrorResponse(c, http.StatusUnauthorized, "User not valid!", nil)
			}

			// Add claims to context
			c.Set("userID", claims.UserID)
			c.Set("userRole", claims.Role)

			return next(c)
		}
	}
}
