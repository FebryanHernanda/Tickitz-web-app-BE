package repositories

import (
	"context"

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

func (r *SeatRepository) GetAvailableSeats(ctx context.Context, cinemaScheduleID int) ([]models.Seat, error) {
	query := `
	SELECT s.id AS seat_id, s.seat_number, s.seat_type
    FROM seats s
    WHERE s.id NOT IN (
		SELECT os.seat_id
		FROM orders_seats os
		JOIN orders o ON os.order_id = o.id
		WHERE o.cinemas_schedule_id = $1
			AND os.status = 'booked'
        )
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
