package controllers

import (
	"net/http"

	"github.com/dfanso/parkme-backend/internal/dto"
	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/auth"
	"github.com/dfanso/parkme-backend/pkg/s3"
	"github.com/dfanso/parkme-backend/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthController struct {
	userService *services.UserService
	jwtManager  *auth.JWTManager
	s3Client    *s3.S3Client
}

func NewAuthController(userService *services.UserService, jwtManager *auth.JWTManager, s3Client *s3.S3Client) *AuthController {
	return &AuthController{
		userService: userService,
		jwtManager:  jwtManager,
		s3Client:    s3Client,
	}
}

// GetJWTManager returns the JWT manager instance
func (c *AuthController) GetJWTManager() *auth.JWTManager {
	return c.jwtManager
}

func (c *AuthController) Login(ctx echo.Context) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&credentials); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Find user by email using FindOne
	user, err := c.userService.FindOne(ctx.Request().Context(), bson.M{"email": credentials.Email})
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Compare password
	if err := user.ComparePassword(credentials.Password); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// Generate JWT token with correct arguments
	token, err := c.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) Register(ctx echo.Context) error {
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err)
	}

	// Find user by email
	filter := bson.M{"email": req.Email}
	existingUser, err := c.userService.FindOne(ctx.Request().Context(), filter)
	if err == nil && existingUser != nil {
		return utils.ErrorResponse(ctx, http.StatusConflict, "user already exists", nil)
	}

	// Create new user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     models.RoleUser,
		Status:   models.UserStatusActive,
	}
	if err := user.BeforeCreate(); err != nil {
		// Handle validation errors specifically
		if e, ok := err.(validation.Errors); ok {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", e)
		}
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user data", err)
	}

	if err := c.userService.Create(ctx.Request().Context(), user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create user", err)
	}

	// Generate JWT token
	token, err := c.jwtManager.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate token", err)
	}
	user.Password = "" // Clear password

	return utils.SuccessResponse(ctx, http.StatusOK, "Login successful", dto.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile retrieves the user profile using the JWT token
func (c *AuthController) GetProfile(ctx echo.Context) error {
	// Get userID from context (set by AuthMiddleware)
	userIDInterface := ctx.Get("userID")
	if userIDInterface == nil {
		return utils.ErrorResponse(ctx, http.StatusUnauthorized, "User ID not found in token", nil)
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID format", nil)
	}

	// Get user profile from database
	user, err := c.userService.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", err)
	}

	// Clear sensitive information
	user.Password = ""

	return utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", user)
}

func (c *AuthController) UpdateProfile(ctx echo.Context) error {
	// Get userID from context
	userID := ctx.Get("userID").(primitive.ObjectID)

	// Get current user
	user, err := c.userService.GetByID(ctx.Request().Context(), userID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", err)
	}

	// Update name if provided
	name := ctx.FormValue("name")
	if name != "" {
		user.Name = name
	}

	// Handle file upload if provided
	file, err := ctx.FormFile("image")
	if err == nil { // File was provided
		s3URL, err := utils.UploadFileToS3(file, c.s3Client, "profile-images")
		if err != nil {
			return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to upload image", err)
		}
		user.ProfileImageURL = s3URL
	}

	// Update user in database
	if err := c.userService.Update(ctx.Request().Context(), user); err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update profile", err)
	}

	user.Password = "" // Clear password before sending response
	return utils.SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", user)
}
