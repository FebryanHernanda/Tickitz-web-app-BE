package models

type Seat struct {
	SeatID     int    `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`
}
