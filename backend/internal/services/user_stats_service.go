package services

import (
	"context"
	"time"

	"github.com/dfanso/parkme-backend/internal/dto"
	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/dfanso/parkme-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStatsService struct {
	bookingRepo *repositories.BookingRepository
}

func NewUserStatsService(bookingRepo *repositories.BookingRepository) *UserStatsService {
	return &UserStatsService{
		bookingRepo: bookingRepo,
	}
}

func (s *UserStatsService) GetUserStats(ctx context.Context, userID primitive.ObjectID) (*dto.UserStatsResponse, error) {
	// Get all user bookings
	bookings, err := s.bookingRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &dto.UserStatsResponse{}
	var totalDuration float64
	now := time.Now()

	// Calculate statistics
	for _, booking := range bookings {
		stats.TotalBookings++

		if booking.EndTime != nil {
			duration := booking.EndTime.Sub(booking.StartTime).Hours()
			totalDuration += duration
		}

		// Calculate booking status
		switch booking.Status {
		case models.BookingStatusCancelled:
			stats.CancelledBookings++
		case models.BookingStatusCompleted:
			stats.CompletedBookings++
			if booking.TotalAmount != nil {
				stats.TotalSpentAmount += *booking.TotalAmount
			}
		case models.BookingStatusActive, models.BookingStatusPending:
			stats.ActiveBookings++
			if stats.UpcomingBookingDate == "" && booking.StartTime.After(now) {
				stats.UpcomingBookingDate = booking.StartTime.Format(time.RFC3339)
			}
		}

		// Track last booking (using CreatedAt for consistent ordering)
		if stats.LastBookingDate == "" || booking.CreatedAt.After(now) {
			stats.LastBookingDate = booking.StartTime.Format(time.RFC3339)
		}
	}

	// Calculate average duration if there are completed bookings
	if stats.CompletedBookings > 0 {
		stats.AverageBookingDuration = totalDuration / float64(stats.CompletedBookings)
	}

	return stats, nil
}
