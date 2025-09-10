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

func (pr *ProfileRepository) UpdateProfile(ctx context.Context, userID int, update *models.UserUpdateRequest) error {
	profileSet := []string{}
	profileArgs := []any{}
	argPos := 1

	if update.Profile.FirstName != nil {
		profileSet = append(profileSet, fmt.Sprintf("first_name = $%d", argPos))
		profileArgs = append(profileArgs, *update.Profile.FirstName)
		argPos++
	}
	if update.Profile.LastName != nil {
		profileSet = append(profileSet, fmt.Sprintf("last_name = $%d", argPos))
		profileArgs = append(profileArgs, *update.Profile.LastName)
		argPos++
	}
	if update.Profile.PhoneNumber != nil {
		profileSet = append(profileSet, fmt.Sprintf("phone_number = $%d", argPos))
		profileArgs = append(profileArgs, *update.Profile.PhoneNumber)
		argPos++
	}
	if update.Profile.ImagePath != nil {
		profileSet = append(profileSet, fmt.Sprintf("image_path = $%d", argPos))
		profileArgs = append(profileArgs, *update.Profile.ImagePath)
		argPos++
	}
	if update.Profile.Points != nil {
		profileSet = append(profileSet, fmt.Sprintf("points = $%d", argPos))
		profileArgs = append(profileArgs, *update.Profile.Points)
		argPos++
	}

	if len(profileSet) > 0 {
		query := fmt.Sprintf("UPDATE profiles SET %s WHERE user_id = $%d", strings.Join(profileSet, ", "), argPos)
		profileArgs = append(profileArgs, userID)

		if _, err := pr.DB.Exec(ctx, query, profileArgs...); err != nil {
			return err
		}
	}

	userSet := []string{}
	userArgs := []any{}
	userPos := 1

	if update.User.Email != nil {
		userSet = append(userSet, fmt.Sprintf("email = $%d", userPos))
		userArgs = append(userArgs, update.User.Email)
		userPos++
	}
	if update.User.Password != nil {
		userSet = append(userSet, fmt.Sprintf("password = $%d", userPos))
		userArgs = append(userArgs, update.User.Password)
		userPos++
	}

	if len(userSet) > 0 {
		userSet = append(userSet, "updated_at = NOW()")
		query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(userSet, ", "), userPos)
		userArgs = append(userArgs, userID)

		if _, err := pr.DB.Exec(ctx, query, userArgs...); err != nil {
			return err
		}
	} else {
		if _, err := pr.DB.Exec(ctx, "UPDATE users SET updated_at = NOW() WHERE id = $1", userID); err != nil {
			return err
		}
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
