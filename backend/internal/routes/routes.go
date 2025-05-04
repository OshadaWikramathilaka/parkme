package routes

import (
	"github.com/dfanso/parkme-backend/internal/controllers"
	customMiddleware "github.com/dfanso/parkme-backend/pkg/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(
	e *echo.Echo,
	userController *controllers.UserController,
	authController *controllers.AuthController,
	vehicleController *controllers.VehicleController,
	arduinoController *controllers.ArduinoController,
	bookingController *controllers.BookingController,
	walletController *controllers.WalletController,
	parkingLocationController *controllers.ParkingLocationController,
	userStatsController *controllers.UserStatsController,
) {
	// Public routes
	auth := e.Group("/api/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	// Arduino routes
	arduino := e.Group("/api/arduino")
	arduino.Use(customMiddleware.ValidateAPIKey())
	arduino.POST("/gate/enter/upload", arduinoController.GateEnter)
	arduino.POST("/gate/exit/upload", arduinoController.GateExit)
	arduino.POST("/spot/status", arduinoController.UpdateSpotStatus)

	// Protected routes
	api := e.Group("/api")
	api.Use(customMiddleware.AuthMiddleware(authController.GetJWTManager(), userController.GetUserService()))

	// User routes
	users := api.Group("/users")
	users.GET("", userController.GetAll)
	users.GET("/:id", userController.GetByID)
	users.PUT("/:id", userController.Update)
	users.DELETE("/:id", userController.Delete)

	// User Stats route
	api.GET("/user/stats", userStatsController.GetUserStats)

	// Auth routes for protected endpoints
	auth = api.Group("/auth")
	auth.GET("/profile", authController.GetProfile)
	auth.PUT("/profile", authController.UpdateProfile)

	// Vehicle routes
	vehicles := api.Group("/vehicles")
	vehicles.POST("", vehicleController.Create)
	vehicles.GET("", vehicleController.GetAll)
	vehicles.GET("/:id", vehicleController.GetByID)
	vehicles.GET("/user/:userId", vehicleController.GetUserVehicles)
	vehicles.PUT("/:id", vehicleController.Update)
	vehicles.DELETE("/:id", vehicleController.Delete)

	// Booking routes
	bookings := api.Group("/bookings")
	bookings.POST("", bookingController.CreateBooking)
	bookings.GET("/:id", bookingController.GetBooking)
	bookings.GET("/user", bookingController.GetUserBookings)
	bookings.GET("/vehicle/:id", bookingController.GetVehicleBookings)
	bookings.PUT("/:id/cancel", bookingController.CancelBooking)

	// Wallet routes
	wallet := api.Group("/wallet")
	wallet.POST("/topup", walletController.TopUp)
	wallet.GET("/balance", walletController.GetBalance)
	wallet.GET("/transactions", walletController.GetTransactions)

	// Parking Location routes
	locations := api.Group("/locations")
	locations.POST("", parkingLocationController.CreateLocation)
	locations.GET("", parkingLocationController.GetAllLocations)
	locations.GET("/:id", parkingLocationController.GetLocation)
	locations.PUT("/:id", parkingLocationController.UpdateLocation)
	locations.PUT("/:id/slot", parkingLocationController.UpdateSlotStatus)
	locations.DELETE("/:id", parkingLocationController.DeleteLocation)
}
