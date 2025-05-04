package dto

type UserStatsResponse struct {
	TotalBookings          int      `json:"total_bookings"`
	ActiveBookings         int      `json:"active_bookings"`
	CompletedBookings      int      `json:"completed_bookings"`
	CancelledBookings      int      `json:"cancelled_bookings"`
	TotalSpentAmount       float64  `json:"total_spent_amount"`
	AverageBookingDuration float64  `json:"average_booking_duration"` // in hours
	LastBookingDate        string   `json:"last_booking_date,omitempty"`
	UpcomingBookingDate    string   `json:"upcoming_booking_date,omitempty"`
	FavoriteLocations      []string `json:"favorite_locations,omitempty"`
}
