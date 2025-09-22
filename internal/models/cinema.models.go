package models

import "time"

type CinemaSeat struct {
	SeatID     int    `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`
}

type GetFilterSchedules struct {
	CinemaScheduleID int
	CinemaName       *string
	TicketPrice      float64
	LocationName     *string    `form:"location_name"`
	ScheduleDate     *time.Time `form:"schedule_date"`
	ScheduleTime     *string    `form:"schedule_time"`
	MovieName        *string    `form:"movie_name"`
}

type CinemaList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CinemaLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetFilterSchedulesCache struct {
	Data       []GetFilterSchedules `json:"data"`
	TotalCount int                  `json:"total_count"`
}
