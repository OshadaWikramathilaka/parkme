package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrBookingNotFound          = errors.New("booking not found")
	ErrSpotAlreadyBooked        = errors.New("parking spot already booked for this time")
	ErrInvalidTimeRange         = errors.New("invalid time range")
	ErrBookingInPast            = errors.New("cannot book in the past")
	ErrInsufficientWalletPoints = errors.New("insufficient wallet balance, minimum 300 points required")
	ErrNoAvailableSlots         = errors.New("no available parking slots at this location")
)

type BookingService struct {
	repo     *repositories.BookingRepository
	vehicle  *VehicleService
	wallet   *WalletService
	location *ParkingLocationService
	user     *UserService
}

func NewBookingService(repo *repositories.BookingRepository, vehicleService *VehicleService, walletService *WalletService, locationService *ParkingLocationService, userService *UserService) *BookingService {
	return &BookingService{
		repo:     repo,
		vehicle:  vehicleService,
		wallet:   walletService,
		location: locationService,
		user:     userService,
	}
}

// findAvailableSlot finds an available parking slot at the given location
func (s *BookingService) findAvailableSlot(ctx context.Context, locationID primitive.ObjectID) (string, error) {
	location, err := s.location.GetLocation(ctx, locationID)
	if err != nil {
		return "", err
	}

	for _, slot := range location.Slots {
		if !slot.IsOccupied {
			return slot.Number, nil
		}
	}

	return "", ErrNoAvailableSlots
}

func (s *BookingService) CreateBooking(ctx context.Context, booking *models.Booking) error {
	// Check wallet balance first
	balance, err := s.wallet.GetBalance(ctx, booking.UserID)
	if err != nil {
		return err
	}

	if balance < 300 {
		return ErrInsufficientWalletPoints
	}

	// Set booking type based on whether it's pre-booked or on-site
	if booking.SpotNumber != nil {
		booking.BookingType = models.BookingTypePreBooked
	} else {
		booking.BookingType = models.BookingTypeOnSite
	}

	// For on-site booking, we assume it's for the current time, so skip the past booking check
	// Only check for pre-booked bookings
	if booking.BookingType == models.BookingTypePreBooked && booking.StartTime.Before(time.Now()) {
		return ErrBookingInPast
	}

	// Debug logging with proper pointer handling
	spotNumberStr := "nil"
	if booking.SpotNumber != nil {
		spotNumberStr = *booking.SpotNumber
	}
	fmt.Printf("Processing booking with spot number: %s, type: %s\n", spotNumberStr, booking.BookingType)

	// Verify location exists and has available slots
	if booking.BookingType == models.BookingTypePreBooked {
		// Validate required fields for pre-booked
		if booking.SpotNumber == nil {
			return errors.New("spot number is required for pre-booked parking")
		}

		// If end time is provided, validate it
		if booking.EndTime != nil && booking.StartTime.After(*booking.EndTime) {
			return ErrInvalidTimeRange
		}

		// Check if the spot exists in the location
		location, err := s.location.GetLocation(ctx, booking.LocationID)
		if err != nil {
			return fmt.Errorf("failed to get location: %w", err)
		}

		spotExists := false
		for _, slot := range location.Slots {
			if slot.Number == *booking.SpotNumber {
				spotExists = true
				break
			}
		}

		if !spotExists {
			return fmt.Errorf("spot number %s does not exist in this location", *booking.SpotNumber)
		}

		// If end time is provided, check availability for the time range
		if booking.EndTime != nil {
			if !s.isSlotAvailable(ctx, *booking.SpotNumber, booking.StartTime, *booking.EndTime) {
				return ErrSpotAlreadyBooked
			}
		} else {
			// If no end time, just check if the spot is currently available
			if !s.isSpotCurrentlyAvailable(ctx, *booking.SpotNumber) {
				return ErrSpotAlreadyBooked
			}
		}
	} else {
		// For on-site booking, find an available slot
		spotNumber, err := s.findAvailableSlot(ctx, booking.LocationID)
		if err != nil {
			return err
		}
		booking.SpotNumber = &spotNumber
		fmt.Printf("Assigned spot number %s for on-site booking\n", spotNumber)

		// Update slot status to occupied
		err = s.location.UpdateSlotStatus(ctx, booking.LocationID, spotNumber, true)
		if err != nil {
			return err
		}

		// For on-site booking, set start time to now if not provided
		if booking.StartTime.IsZero() {
			booking.StartTime = time.Now()
		}
	}

	// Verify vehicle exists
	_, err = s.vehicle.GetByID(ctx, booking.VehicleID)
	if err != nil {
		return err
	}

	// Set timestamps and status
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()
	booking.Status = models.BookingStatusPending

	return s.repo.Create(ctx, booking)
}

func (s *BookingService) GetBooking(ctx context.Context, id primitive.ObjectID) (*models.Booking, error) {
	booking, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return booking, nil
}

func (s *BookingService) GetVehicleBookings(ctx context.Context, vehicleID primitive.ObjectID) ([]models.Booking, error) {
	return s.repo.FindByVehicle(ctx, vehicleID)
}

func (s *BookingService) UpdateBookingStatus(ctx context.Context, id primitive.ObjectID, status models.BookingStatus) error {
	booking, err := s.GetBooking(ctx, id)
	if err != nil {
		return err
	}

	booking.Status = status
	booking.UpdatedAt = time.Now()

	return s.repo.Update(ctx, booking)
}

func (s *BookingService) CancelBooking(ctx context.Context, id primitive.ObjectID) error {
	//only cancel if the booking is pending
	booking, err := s.GetBooking(ctx, id)
	if err != nil {
		return err
	}
	if booking.Status != models.BookingStatusPending {
		return errors.New("booking is not pending")
	}
	return s.UpdateBookingStatus(ctx, id, models.BookingStatusCancelled)
}

func (s *BookingService) isSlotAvailable(ctx context.Context, spotNumber string, start, end time.Time) bool {
	filter := bson.M{
		"spot_number": spotNumber,
		"status": bson.M{
			"$in": []models.BookingStatus{
				models.BookingStatusActive,
				models.BookingStatusPending,
			},
		},
		"$or": []bson.M{
			{
				"start_time": bson.M{
					"$lt": end,
				},
				"end_time": bson.M{
					"$gt": start,
				},
			},
			{
				"start_time": bson.M{
					"$eq": start,
				},
			},
		},
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		fmt.Printf("Error checking slot availability for time range: %v\n", err)
		return false
	}

	isAvailable := count == 0
	fmt.Printf("Checking time range availability for spot %s: start=%v, end=%v, count=%d, available=%v\n",
		spotNumber, start, end, count, isAvailable)
	return isAvailable
}

// FindBookingByFilter finds a booking using the provided filter
func (s *BookingService) FindBookingByFilter(ctx context.Context, filter bson.M) (*models.Booking, error) {
	booking, err := s.repo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return booking, nil
}

// Add method to complete on-site booking
func (s *BookingService) CompleteOnSiteBooking(ctx context.Context, bookingID primitive.ObjectID, endTime time.Time, totalAmount float64) error {
	booking, err := s.repo.FindByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking.BookingType != models.BookingTypeOnSite {
		return errors.New("only on-site bookings can be completed this way")
	}

	booking.EndTime = &endTime
	booking.TotalAmount = &totalAmount
	booking.Status = models.BookingStatusCompleted
	booking.UpdatedAt = time.Now()

	return s.repo.Update(ctx, booking)
}

func (s *BookingService) GetUserBookings(ctx context.Context, userID primitive.ObjectID) ([]models.Booking, error) {
	//populate user
	bookings, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i, booking := range bookings {
		user, err := s.user.GetByID(ctx, booking.UserID)
		if err != nil {
			return nil, err
		}
		bookings[i].User = user
	}

	//populate vehicle
	for i, booking := range bookings {
		vehicle, err := s.vehicle.GetByID(ctx, booking.VehicleID)
		if err != nil {
			return nil, err
		}
		bookings[i].Vehicle = vehicle
	}

	//populate location
	for i, booking := range bookings {
		location, err := s.location.GetLocation(ctx, booking.LocationID)
		if err != nil {
			return nil, err
		}
		bookings[i].Location = location
	}
	return bookings, nil
}

// isSpotCurrentlyAvailable checks if a spot is currently available (no active or pending bookings)
func (s *BookingService) isSpotCurrentlyAvailable(ctx context.Context, spotNumber string) bool {
	now := time.Now()
	filter := bson.M{
		"spot_number": spotNumber,
		"status": bson.M{
			"$in": []models.BookingStatus{
				models.BookingStatusActive,
				models.BookingStatusPending,
			},
		},
		"start_time": bson.M{
			"$lte": now,
		},
		"$or": []bson.M{
			{
				"end_time": nil,
			},
			{
				"end_time": bson.M{
					"$gt": now,
				},
			},
		},
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		fmt.Printf("Error checking spot availability: %v\n", err)
		return false
	}

	isAvailable := count == 0
	fmt.Printf("Checking current availability for spot %s: count=%d, available=%v\n", spotNumber, count, isAvailable)
	return isAvailable
}
