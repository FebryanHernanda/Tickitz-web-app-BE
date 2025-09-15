package repositories

import (
	"context"
	"log"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CinemaRepository struct {
	DB *pgxpool.Pool
}

func NewCinemaRepository(db *pgxpool.Pool) *CinemaRepository {
	return &CinemaRepository{
		DB: db,
	}
}

func (r *CinemaRepository) IsCinemaScheduleExists(ctx context.Context, cinemaScheduleID int) (bool, error) {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM cinemas_schedules WHERE id = $1)`
	err := r.DB.QueryRow(ctx, query, cinemaScheduleID).Scan(&exist)
	if err != nil {
		log.Printf("ERROR \nCause :  %s", err)
		return false, err
	}
	return exist, nil
}

func (r *CinemaRepository) GetAvailableSeats(ctx context.Context, cinemaScheduleID int) ([]models.CinemaSeat, error) {

	query := `
	SELECT
		os.seat_id,
		s.seat_number,
		s.seat_type
	FROM
		orders_seats os
		JOIN seats s ON os.seat_id = s.id
		JOIN orders o ON os.order_id = o.id
	WHERE
		o.cinemas_schedule_id = $1
		AND os.status = 'booked'
	`

	rows, err := r.DB.Query(ctx, query, cinemaScheduleID)
	if err != nil {
		return nil, err
	}

	var seats []models.CinemaSeat
	for rows.Next() {
		var s models.CinemaSeat
		err := rows.Scan(
			&s.SeatID,
			&s.SeatNumber,
			&s.SeatType,
		)
		if err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil
}

func (r *CinemaRepository) GetScheduleFilter(ctx context.Context, locationFilter *string, dateFilter *time.Time, timeFilter *string, limit, offset int) ([]models.GetFilterSchedules, error) {
	query := `
	SELECT
		c.name AS cinema_name,
		l.name AS location_name,
		s.date AS schedule_date,
		s.time AS schedule_time,
		m.title AS movie_title
	FROM
    	cinemas_schedules cs
	JOIN
    	cinemas c ON cs.cinemas_id = c.id
	JOIN
    	locations l ON cs.locations_id = l.id
	JOIN
    	schedules s ON cs.schedules_id = s.id
	JOIN 
		movies m ON s.movie_id = m.id
	WHERE
		($1::text IS NULL OR l.name = $1::text)
		AND ($2::date IS NULL OR s.date = $2::date)
		AND ($3::show_time IS NULL OR s.time = $3::show_time)
	ORDER BY
    	s.time ASC
	LIMIT $4 OFFSET $5
`

	values := []any{locationFilter, dateFilter, timeFilter, limit, offset}

	rows, err := r.DB.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var filterSchedule []models.GetFilterSchedules
	for rows.Next() {
		var fs models.GetFilterSchedules
		err := rows.Scan(
			&fs.CinemaName,
			&fs.LocationName,
			&fs.ScheduleDate,
			&fs.ScheduleTime,
			&fs.MovieName,
		)
		if err != nil {
			return nil, err
		}
		filterSchedule = append(filterSchedule, fs)
	}

	return filterSchedule, nil
}
