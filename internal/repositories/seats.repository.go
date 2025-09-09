package repositories

import (
	"context"
	"log"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	DB *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{
		DB: db,
	}
}

func (r *SeatRepository) IsCinemaScheduleExists(ctx context.Context, cinemaScheduleID int) (bool, error) {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM cinemas_schedules WHERE id = $1)`
	err := r.DB.QueryRow(ctx, query, cinemaScheduleID).Scan(&exist)
	if err != nil {
		log.Printf("ERROR \nCause :  %s", err)
		return false, err
	}
	return exist, nil
}

func (r *SeatRepository) GetAvailableSeats(ctx context.Context, cinemaScheduleID int) ([]models.Seat, error) {

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

	var seats []models.Seat
	for rows.Next() {
		var s models.Seat
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
