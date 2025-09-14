package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

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

// helper
func (r *AdminRepository) IsMoviesExists(ctx context.Context, movieID int) (bool, error) {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)`
	err := r.DB.QueryRow(ctx, query, movieID).Scan(&exist)
	if err != nil {
		log.Printf("ERROR \nCause :  %s", err)
		return false, err
	}

	return exist, nil
}

// helper
func (r *AdminRepository) IsScheduleExists(ctx context.Context, scheduleID int) (bool, error) {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM schedules WHERE id = $1)`
	err := r.DB.QueryRow(ctx, query, scheduleID).Scan(&exist)
	if err != nil {
		log.Printf("ERROR \nCause :  %s", err)
		return false, err
	}

	return exist, nil
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
			COALESCE(ARRAY_AGG(DISTINCT sch.time ORDER BY sch.time ASC) FILTER (WHERE sch.time IS NOT NULL),'{}') AS time_playing,
			COALESCE(ARRAY_AGG(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL),'{}') AS casts,
			COALESCE(ARRAY_AGG(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL),'{}') AS genres
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
	defer dbTx.Rollback(ctx)

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

func (r *AdminRepository) UpdateMovies(ctx context.Context, id int, update models.EditMovies) error {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer dbTx.Rollback(ctx)

	updateData := map[string]interface{}{}
	if update.Title != nil {
		updateData["title"] = *update.Title
	}
	if update.Synopsis != nil {
		updateData["synopsis"] = *update.Synopsis
	}
	if update.ReleaseDate != nil {
		updateData["release_date"] = *update.ReleaseDate
	}
	if update.Rating != nil {
		updateData["rating"] = *update.Rating
	}
	if update.AgeRating != nil {
		updateData["age_rating"] = *update.AgeRating
	}
	if update.Duration != nil {
		updateData["duration"] = *update.Duration
	}
	if update.DirectorID != nil {
		updateData["director_id"] = *update.DirectorID
	}
	if update.PosterPath != nil {
		updateData["poster_path"] = *update.PosterPath
	}
	if update.BackdropPath != nil {
		updateData["backdrop_path"] = *update.BackdropPath
	}

	set := []string{}
	args := []any{}
	i := 1
	for k, v := range updateData {
		set = append(set, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	if len(set) > 0 {
		query := fmt.Sprintf("UPDATE movies SET %s WHERE id = $%d", strings.Join(set, ", "), i)
		args = append(args, id)
		_, err := dbTx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	if update.Genres != nil {
		_, err = dbTx.Exec(ctx, "DELETE FROM movies_genres WHERE movie_id = $1", id)
		if err != nil {
			return err
		}
		for _, g := range *update.Genres {
			_, err := dbTx.Exec(ctx, "INSERT INTO movies_genres (movie_id, genre_id) VALUES ($1, $2)", id, g)
			if err != nil {
				return err
			}
		}
	}

	if update.Casts != nil {
		_, err = dbTx.Exec(ctx, "DELETE FROM movies_cast WHERE movie_id = $1", id)
		if err != nil {
			return err
		}
		for _, c := range *update.Casts {
			_, err := dbTx.Exec(ctx, "INSERT INTO movies_cast (movie_id, cast_id) VALUES ($1, $2)", id, c)
			if err != nil {
				return err
			}
		}
	}

	if update.Schedules != nil {
		_, err = dbTx.Exec(ctx, "DELETE FROM schedules WHERE movie_id = $1", id)
		if err != nil {
			return err
		}
		for _, s := range *update.Schedules {
			_, err := dbTx.Exec(ctx, "INSERT INTO schedules (movie_id, date, time) VALUES ($1, $2, $3)", id, s.Date, s.Time)
			if err != nil {
				return err
			}
		}
	}

	return dbTx.Commit(ctx)
}

func (r *AdminRepository) GetMovieSchedule(ctx context.Context) ([]models.GetSchedule, error) {
	query := `
	SELECT 
		s.id, 
		s.date, s.
		time, 
		s.movie_id, 
		m.title AS movie_title
	FROM schedules s
	JOIN movies m  ON s.movie_id = m.id
	`

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
			&sch.MovieID,
			&sch.MovieName,
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
