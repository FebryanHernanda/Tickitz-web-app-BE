package repositories

import (
	"context"
	"fmt"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	DB *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{
		DB: db,
	}
}

func (r *AdminRepository) GetAllMovies(ctx context.Context) ([]models.AdminMovies, error) {
	query := `
	  SELECT
            m.id, 
			m.title, 
			m.poster_path, 
			m.backdrop_path, 
			m.synopsis,
            m.release_date, 
			m.rating, 
			m.age_rating, 
			m.duration,
            d.name AS director, 
			sch.date AS date_playing, 
			l.name AS location_name,
            cnm.name AS cinema_name,
            ARRAY_AGG(DISTINCT sch.time ORDER BY sch.time ASC) AS time_playing,
            ARRAY_AGG(DISTINCT c.name) AS casts,
            ARRAY_AGG(DISTINCT g.name) AS genres
        FROM movies m
            LEFT JOIN movies_genres mg ON m.id = mg.movie_id
            LEFT JOIN genres g ON mg.genre_id = g.id
            LEFT JOIN directors d ON m.director_id = d.id
            LEFT JOIN movies_cast mc ON m.id = mc.movie_id
            LEFT JOIN casts c ON mc.cast_id = c.id
            LEFT JOIN schedules sch ON m.id = sch.movie_id
            LEFT JOIN cinemas_schedules cs ON sch.id = cs.schedules_id
            LEFT JOIN locations l ON cs.locations_id = l.id
            LEFT JOIN cinemas cnm on cs.cinemas_id = cnm.id
        GROUP BY m.id, d.id, sch.date, l.id, cnm.id
	`

	var allMovies []models.AdminMovies
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var am models.AdminMovies
		err := rows.Scan(
			&am.ID,
			&am.Title,
			&am.PosterPath,
			&am.BackdropPath,
			&am.Synopsis,
			&am.ReleaseDate,
			&am.Rating,
			&am.AgeRating,
			&am.Duration,
			&am.Director,
			&am.DatePlaying,
			&am.LocationName,
			&am.CinemaName,
			&am.TimePlaying,
			&am.Casts,
			&am.Genres,
		)
		if err != nil {
			return nil, err
		}
		allMovies = append(allMovies, am)
	}

	return allMovies, nil
}

func (r *AdminRepository) AddMovies(ctx context.Context, movie *models.AddMovies) (*models.AddMovies, error) {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed begin db transaction : %w", err)
	}
	defer func() {
		if err != nil {
			dbTx.Rollback(ctx)
		}
	}()

	var movieID int

	queryInsert := `
		INSERT INTO movies (title, poster_path, backdrop_path, synopsis, release_date, rating, age_rating, duration, director_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`
	values := []any{movie.Title, movie.PosterPath, movie.BackdropPath, movie.Synopsis, movie.ReleaseDate,
		movie.Rating, movie.AgeRating, movie.Duration, movie.DirectorID}

	err = dbTx.QueryRow(ctx, queryInsert, values...).Scan(&movieID)
	if err != nil {
		return nil, fmt.Errorf("insert movie failed : %w", err)
	}

	queryInsertGenres := `INSERT INTO movies_genres (movie_id, genre_id) VALUES ($1, $2)`
	for _, genreID := range movie.Genres {
		_, err = dbTx.Exec(ctx, queryInsertGenres, movieID, genreID)
		if err != nil {
			return nil, fmt.Errorf("insert movies_genres failed : %w", err)
		}
	}

	queryInsertCast := `INSERT INTO movies_cast (movie_id, cast_id) VALUES ($1, $2)`
	for _, castID := range movie.Casts {
		_, err = dbTx.Exec(ctx, queryInsertCast, movieID, castID)
		if err != nil {
			return nil, fmt.Errorf("insert movies_cast failed : %w", err)
		}
	}

	queryInsertSchedule := `INSERT INTO schedules (movie_id, date, time) VALUES ($1, $2, $3)`
	for _, schedule := range movie.Schedules {
		_, err = dbTx.Exec(ctx, queryInsertSchedule, movieID, schedule.Date, schedule.Time)
		if err != nil {
			return nil, fmt.Errorf("insert schedule failed : %w", err)
		}
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit db transaction failed : %w", err)
	}

	movieData := *movie

	return &movieData, nil
}

func (r *AdminRepository) AddCinemaSchedule(ctx context.Context, data []models.CinemaScheduleLocation) error {
	query := "INSERT INTO cinemas_schedules (cinemas_id, schedules_id, locations_id) VALUES ($1, $2, $3)"
	for _, item := range data {
		_, err := r.DB.Exec(ctx, query, item.CinemaID, item.ScheduleID, item.LocationID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *AdminRepository) GetMovieSchedule(ctx context.Context) ([]models.GetSchedule, error) {
	query := `SELECT id, date, time FROM schedules`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedule []models.GetSchedule
	for rows.Next() {
		var sch models.GetSchedule
		err := rows.Scan(
			&sch.ID,
			&sch.Date,
			&sch.Time,
		)
		if err != nil {
			return nil, err
		}
		schedule = append(schedule, sch)
	}

	return schedule, nil
}

func (r *AdminRepository) DeleteMovies(ctx context.Context, movieID int) error {
	query := `DELETE FROM movies WHERE id = $1`

	movies, err := r.DB.Exec(ctx, query, movieID)
	if err != nil {
		return err
	}

	if movies.RowsAffected() == 0 {
		return fmt.Errorf("movie not found")
	}

	return nil
}
