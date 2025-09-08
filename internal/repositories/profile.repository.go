package repositories

import (
	"context"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	DB *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (pr *ProfileRepository) GetProfile(ctx context.Context, userID int) (*models.Profile, error) {

	query := `
        SELECT
            u.id,
            u.email,
            u.role,
            u.virtual_account,
            p.first_name,
            p.last_name,
            p.phone_number,
            p.points,
            p.image_path
        FROM users u
        LEFT JOIN profiles p ON u.id = p.user_id
        WHERE u.id = $1
    `

	var userProfile models.Profile
	rows := pr.DB.QueryRow(ctx, query, userID)
	err := rows.Scan(
		&userProfile.ID,
		&userProfile.Email,
		&userProfile.Role,
		&userProfile.VirtualAccount,
		&userProfile.FirstName,
		&userProfile.LastName,
		&userProfile.PhoneNumber,
		&userProfile.Points,
		&userProfile.ImagePath,
	)
	if err != nil {
		return nil, err
	}

	return &userProfile, err
}
