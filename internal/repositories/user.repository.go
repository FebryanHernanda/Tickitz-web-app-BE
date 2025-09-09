package repositories

import (
	"context"
	"time"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	queryUser := `
	INSERT INTO users (email, password, role,  virtual_account, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	values := []any{user.Email, user.Password, user.Role, user.VirtualAccount, user.CreatedAt, user.UpdatedAt}

	var userID int
	err := r.DB.QueryRow(ctx, queryUser, values...).Scan(&userID)
	if err != nil {
		return err
	}

	queryProfile := `
        INSERT INTO profiles (user_id, points, image_path)
        VALUES ($1, 0, 'public/profile/default.png')
		`
	_, err = r.DB.Exec(ctx, queryProfile, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, email, password, role FROM users WHERE email=$1`

	err := r.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
