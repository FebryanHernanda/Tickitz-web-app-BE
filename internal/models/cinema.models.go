package models

import "time"

type CinemaSeat struct {
	SeatID     int    `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`
}

type GetFilterSchedules struct {
	CinemaName   *string
	LocationName *string    `form:"location_name"`
	ScheduleDate *time.Time `form:"schedule_date"`
	ScheduleTime *string    `form:"schedule_time"`
}
