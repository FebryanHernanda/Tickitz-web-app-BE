package repositories

import (
	"context"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersRepository struct {
	DB *pgxpool.Pool
}

func NewOrdersRepository(db *pgxpool.Pool) *OrdersRepository {
	return &OrdersRepository{
		DB: db,
	}
}

func (r *OrdersRepository) CreateOrder(ctx context.Context, order *models.Order) (int, error) {
	query := `INSERT INTO orders (qr_code, isPaid, isActive, total_prices, user_id, cinemas_schedule_id, payment_method_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var orderID int
	values := []any{order.QRCode, order.IsPaid, order.IsActive, order.TotalPrices, order.UserID, order.CinemasScheduleID, order.PaymentMethodID}
	err := r.DB.QueryRow(ctx, query, values...).Scan(&orderID)

	return orderID, err
}

func (r *OrdersRepository) AreSeatsAvailable(ctx context.Context, cinemaScheduleID int, seatIDs []int) (bool, error) {
	query := `
        SELECT 1
        FROM orders_seats os
        JOIN orders o ON os.order_id = o.id
        WHERE o.cinemas_schedule_id = $1
          AND os.seat_id = ANY($2)
          AND os.status = 'booked'
        LIMIT 1
    `
	var exists int
	err := r.DB.QueryRow(ctx, query, cinemaScheduleID, seatIDs).Scan(&exists)
	if err != nil && err != pgx.ErrNoRows {
		return false, err
	}
	return err == pgx.ErrNoRows, nil
}

func (r *OrdersRepository) CreateOrderSeats(ctx context.Context, orderID int, seats []models.OrderSeatInput) error {
	query := `INSERT INTO orders_seats (status, order_id, seat_id) VALUES ($1, $2, $3)`

	for _, seat := range seats {
		values := []any{seat.Status, orderID, seat.SeatID}
		_, err := r.DB.Exec(ctx, query, values...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrdersRepository) GetOrdersHistory(ctx context.Context, userID int) ([]models.OrderHistory, error) {
	query := `
        SELECT
            o.id,
            o.isactive,
            o.ispaid,
            o.qr_code,
            o.total_prices,
            o.user_id,
            u.virtual_account,
            m.title,
            m.age_rating,
            c.name AS cinema,
            c.image_path,
            l.name AS location,
            sch.date::text,
            sch.time::text,
            ARRAY_AGG(s.seat_number) AS seat_number,
            s.seat_type
        FROM
            orders o
            LEFT JOIN cinemas_schedules cs ON o.cinemas_schedule_id = cs.id
            LEFT JOIN cinemas c ON cs.cinemas_id = c.id
            LEFT JOIN locations l ON cs.locations_id = l.id
            LEFT JOIN schedules sch ON cs.schedules_id = sch.id
            LEFT JOIN movies m ON sch.movie_id = m.id
            LEFT JOIN users u ON o.user_id = u.id
            LEFT JOIN orders_seats os ON o.id = os.order_id
            LEFT JOIN seats s ON os.seat_id = s.id
        WHERE
            o.user_id = $1
        GROUP BY
            o.id,
            u.virtual_account,
            m.title,
            m.age_rating,
            c.name,
            c.image_path,
            l.name,
            s.seat_type,
            sch.date,
            sch.time,
            o.isactive,
            o.ispaid,
            o.qr_code,
            o.total_prices
        ORDER BY
            o.created_at ASC;
    `

	var orderHistory []models.OrderHistory
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oh models.OrderHistory
		err := rows.Scan(
			&oh.ID,
			&oh.IsActive,
			&oh.IsPaid,
			&oh.QRCode,
			&oh.TotalPrices,
			&oh.UserID,
			&oh.VirtualAccount,
			&oh.Title,
			&oh.AgeRating,
			&oh.Cinema,
			&oh.CinemaImage,
			&oh.Location,
			&oh.Date,
			&oh.Time,
			&oh.SeatNumbers,
			&oh.SeatType,
		)
		if err != nil {
			return nil, err
		}
		orderHistory = append(orderHistory, oh)
	}
	return orderHistory, nil
}
