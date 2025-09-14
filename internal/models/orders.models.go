package models

import "time"

type Order struct {
	ID                int              `json:"id"`
	QRCode            string           `json:"qr_code"`
	IsPaid            bool             `json:"is_paid"`
	IsActive          bool             `json:"is_active"`
	TotalPrices       float64          `json:"total_prices"`
	UserID            int              `json:"user_id"`
	CinemasScheduleID int              `json:"cinemas_schedule_id"`
	PaymentMethodID   int              `json:"payment_method_id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	OrderSeats        []OrderSeatInput `json:"seats"`
}

type OrderRequest struct {
	IsPaid            bool             `json:"is_paid" `
	IsActive          bool             `json:"is_active"`
	TotalPrices       float64          `json:"total_prices" binding:"required" example:"120000"`
	CinemasScheduleID int              `json:"cinemas_schedule_id" binding:"required" example:"1"`
	PaymentMethodID   int              `json:"payment_method_id" binding:"required" example:"2"`
	OrderSeats        []OrderSeatInput `json:"seats" binding:"required,dive"`
}

type OrderSeat struct {
	ID        int
	Status    string
	OrderID   int
	SeatID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderSeatInput struct {
	Status string `json:"status"`
	SeatID int    `json:"seat_id"`
}

type OrderHistory struct {
	ID             int      `json:"id"`
	IsActive       bool     `json:"is_active"`
	IsPaid         bool     `json:"is_paid"`
	QRCode         string   `json:"qr_code"`
	TotalPrices    float64  `json:"total_prices"`
	UserID         int      `json:"user_id"`
	VirtualAccount string   `json:"virtual_account"`
	Title          string   `json:"title"`
	AgeRating      string   `json:"age_rating"`
	Cinema         string   `json:"cinema"`
	CinemaImage    string   `json:"cinema_image"`
	Location       string   `json:"location"`
	Date           string   `json:"date"`
	Time           string   `json:"time"`
	SeatNumbers    []string `json:"seat_numbers"`
	SeatType       string   `json:"seat_type"`
}
