package main

import (
	"log"

	"github.com/dfanso/parkme-backend/config"
	"github.com/dfanso/parkme-backend/internal/controllers"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"github.com/dfanso/parkme-backend/internal/routes"
	"github.com/dfanso/parkme-backend/internal/services"
	"github.com/dfanso/parkme-backend/pkg/auth"
	"github.com/dfanso/parkme-backend/pkg/database"
	"github.com/dfanso/parkme-backend/pkg/s3"

	customMiddleware "github.com/dfanso/parkme-backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize MongoDB
	db, err := database.NewMongoClient(cfg.MongoDB.URI, cfg.MongoDB.NAME)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = false // Show the Echo banner
	e.HidePort = false   // Show the port number

	// Initialize JWT manager
	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	// Initialize S3 client
	s3Client, err := s3.NewS3Client(cfg.S3BucketName)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Middleware
	e.Use(customMiddleware.NewCustomLogger().Middleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	parkingLocationRepo := repositories.NewParkingLocationRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	vehicleService := services.NewVehicleService(vehicleRepo)
	walletService := services.NewWalletService(walletRepo)
	parkingLocationService := services.NewParkingLocationService(parkingLocationRepo)
	bookingService := services.NewBookingService(bookingRepo, vehicleService, walletService, parkingLocationService, userService)
	arduinoService, err := services.NewArduinoService(cfg, vehicleService)
	if err != nil {
		log.Fatalf("Failed to initialize Arduino service: %v", err)
	}
	userStatsService := services.NewUserStatsService(bookingRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(userService, jwtManager, s3Client)
	userController := controllers.NewUserController(userService)
	vehicleController := controllers.NewVehicleController(vehicleService)
	bookingController := controllers.NewBookingController(bookingService)
	arduinoController := controllers.NewArduinoController(arduinoService, bookingService, userService, parkingLocationService, walletService)
	walletController := controllers.NewWalletController(walletService)
	parkingLocationController := controllers.NewParkingLocationController(parkingLocationService)
	userStatsController := controllers.NewUserStatsController(userStatsService)

	// Register routes
	routes.RegisterRoutes(e, userController, authController, vehicleController, arduinoController, bookingController, walletController, parkingLocationController, userStatsController)

	// Protected routes group
	protected := e.Group("/api")
	protected.Use(customMiddleware.AuthMiddleware(jwtManager, userService))

	// health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))
}
