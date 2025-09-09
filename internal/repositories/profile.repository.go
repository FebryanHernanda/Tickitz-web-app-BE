package repositories

import (
	"context"
	"fmt"
	"strings"

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

func (pr *ProfileRepository) UpdateProfile(ctx context.Context, userID int, update *models.ProfileUpdate) error {
	setParts := []string{}
	args := []any{}
	argPos := 1

	if update.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argPos))
		args = append(args, *update.FirstName)
		argPos++
	}

	if update.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name = $%d", argPos))
		args = append(args, *update.LastName)
		argPos++
	}

	if update.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argPos))
		args = append(args, *update.PhoneNumber)
		argPos++
	}

	if update.ImagePath != nil {
		setParts = append(setParts, fmt.Sprintf("image_path = $%d", argPos))
		args = append(args, *update.ImagePath)
		argPos++
	}

	if update.Points != nil {
		setParts = append(setParts, fmt.Sprintf("points = $%d", argPos))
		args = append(args, *update.Points)
		argPos++
	}

	if len(setParts) == 0 {
		return nil
	}

	// Update profile
	query := fmt.Sprintf("UPDATE profile SET %s WHERE user_id = $%d", strings.Join(setParts, ", "), argPos)
	args = append(args, userID)

	if _, err := pr.DB.Exec(ctx, query, args...); err != nil {
		return err
	}

	// Update users.updated_at
	if _, err := pr.DB.Exec(ctx, "UPDATE users SET updated_at = NOW() WHERE id = $1", userID); err != nil {
		return err
	}

	return nil
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
