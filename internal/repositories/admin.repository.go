package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5"
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

func (r *AdminRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]models.AdminMovies, int, error) {
	var totalCount int
	err := r.DB.QueryRow(ctx, "SELECT COUNT(*) FROM movies").Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

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
		COALESCE(ARRAY_AGG(DISTINCT c.name) FILTER (WHERE c.name IS NOT NULL),'{}') AS casts,
		COALESCE(ARRAY_AGG(DISTINCT g.name) FILTER (WHERE g.name IS NOT NULL),'{}') AS genres
	FROM movies m
	LEFT JOIN movies_genres mg ON m.id = mg.movie_id
	LEFT JOIN genres g ON mg.genre_id = g.id
	LEFT JOIN directors d ON m.director_id = d.id
	LEFT JOIN movies_cast mc ON m.id = mc.movie_id
	LEFT JOIN casts c ON mc.cast_id = c.id
	GROUP BY m.id, d.id
	LIMIT $1 OFFSET $2;
	`

	var allMovies []models.AdminMovies
	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
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
			&am.Casts,
			&am.Genres,
		)
		if err != nil {
			return nil, 0, err
		}

		allMovies = append(allMovies, am)
	}

	return allMovies, totalCount, nil
}

func (r *AdminRepository) GetMovieEditDetail(ctx context.Context, movieID int64) (*models.MovieEditDetail, error) {
	query := `
	SELECT
		m.id,
		m.title,
		m.synopsis,
		m.release_date,
		m.duration,
		m.rating,
		m.age_rating,
		m.poster_path,
		m.backdrop_path,
		d.id AS director_id,
		d.name AS director_name,
		COALESCE(json_agg(DISTINCT g.id) FILTER (WHERE g.id IS NOT NULL), '[]') AS genre_ids,
		COALESCE(json_agg(DISTINCT g.name) FILTER (WHERE g.id IS NOT NULL), '[]') AS genre_names,
		COALESCE(json_agg(DISTINCT c.id) FILTER (WHERE c.id IS NOT NULL), '[]') AS cast_ids,
		COALESCE(json_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '[]') AS cast_names,
		COALESCE(
			(
				SELECT json_agg(row_to_json(t))
				FROM (
					SELECT s.date, array_agg(s.time ORDER BY s.time) AS times
					FROM schedules s
					WHERE s.movie_id = m.id
					GROUP BY s.date
				) t
			),
			'[]'
		) AS schedules,
	    COALESCE(
        (
            SELECT json_agg(row_to_json(cs))
            FROM cinemas_schedules cs
            JOIN schedules s ON s.id = cs.schedules_id
            WHERE s.movie_id = m.id
        ),
        '[]'
    ) AS cinemas_schedules
	FROM movies m
	LEFT JOIN directors d ON d.id = m.director_id
	LEFT JOIN movies_genres mg ON mg.movie_id = m.id
	LEFT JOIN genres g ON g.id = mg.genre_id
	LEFT JOIN movies_cast mc ON mc.movie_id = m.id
	LEFT JOIN casts c ON c.id = mc.cast_id
	LEFT JOIN schedules s ON s.movie_id = m.id
	WHERE m.id = $1
	GROUP BY m.id, d.id;
	`

	row := r.DB.QueryRow(ctx, query, movieID)

	var mv models.MovieEditDetail
	var genreIDsRaw, genreNamesRaw, castIDsRaw, castNamesRaw, schedulesRaw, cinemaSchedulesRaw []byte

	err := row.Scan(
		&mv.ID,
		&mv.Title,
		&mv.Synopsis,
		&mv.ReleaseDate,
		&mv.Duration,
		&mv.Rating,
		&mv.AgeRating,
		&mv.PosterPath,
		&mv.BackdropPath,
		&mv.DirectorID,
		&mv.DirectorName,
		&genreIDsRaw,
		&genreNamesRaw,
		&castIDsRaw,
		&castNamesRaw,
		&schedulesRaw,
		&cinemaSchedulesRaw,
	)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON arrays
	if err := json.Unmarshal(genreIDsRaw, &mv.GenreIDs); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(genreNamesRaw, &mv.GenreNames); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(castIDsRaw, &mv.CastIDs); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(castNamesRaw, &mv.CastNames); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(schedulesRaw, &mv.Schedules); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cinemaSchedulesRaw, &mv.CinemaSchedules); err != nil {
		return nil, err
	}

	return &mv, nil
}

func (r *AdminRepository) AddMovies(ctx context.Context, movie *models.AddMovies) (*models.AddMovies, error) {
	dbTx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin db transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	var movieID int

	// 1. Insert movie
	queryInsertMovie := `
        INSERT INTO movies (title, poster_path, backdrop_path, synopsis, release_date, rating, age_rating, duration, director_id)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
        RETURNING id
    `
	values := []any{
		movie.Title,
		movie.PosterPath,
		movie.BackdropPath,
		movie.Synopsis,
		movie.ReleaseDate,
		movie.Rating,
		movie.AgeRating,
		movie.Duration,
		movie.DirectorID,
	}

	err = dbTx.QueryRow(ctx, queryInsertMovie, values...).Scan(&movieID)
	if err != nil {
		return nil, fmt.Errorf("insert movie failed: %w", err)
	}

	// 2. Insert genres
	queryInsertGenres := `INSERT INTO movies_genres (movie_id, genre_id) VALUES ($1, $2)`
	for _, genreID := range movie.Genres {
		_, err = dbTx.Exec(ctx, queryInsertGenres, movieID, genreID)
		if err != nil {
			return nil, fmt.Errorf("insert movies_genres failed: %w", err)
		}
	}

	// 3. Insert cast
	queryInsertCast := `INSERT INTO movies_cast (movie_id, cast_id) VALUES ($1, $2)`
	for _, castID := range movie.Casts {
		_, err = dbTx.Exec(ctx, queryInsertCast, movieID, castID)
		if err != nil {
			return nil, fmt.Errorf("insert movies_cast failed: %w", err)
		}
	}

	// 4. Insert schedules & get schedule IDs
	queryInsertSchedule := `INSERT INTO schedules (movie_id, date, time) VALUES ($1, $2, $3) RETURNING id`
	scheduleIDs := []int{}
	for _, s := range movie.Schedules {
		var scheduleID int
		err = dbTx.QueryRow(ctx, queryInsertSchedule, movieID, s.Date, s.Time).Scan(&scheduleID)
		if err != nil {
			return nil, fmt.Errorf("insert schedule failed: %w", err)
		}
		scheduleIDs = append(scheduleIDs, scheduleID)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit db transaction failed: %w", err)
	}
	// Return data
	movieData := *movie
	movieData.ID = movieID
	movieData.ScheduleIDs = scheduleIDs

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

	// Update movie main data
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

	if len(updateData) > 0 {
		set := []string{}
		args := []any{}
		i := 1
		for k, v := range updateData {
			set = append(set, fmt.Sprintf("%s = $%d", k, i))
			args = append(args, v)
			i++
		}
		query := fmt.Sprintf("UPDATE movies SET %s WHERE id = $%d", strings.Join(set, ", "), i)
		args = append(args, id)
		if _, err := dbTx.Exec(ctx, query, args...); err != nil {
			return err
		}
	}

	// Update genres
	if update.Genres != nil {
		_, err := dbTx.Exec(ctx, "DELETE FROM movies_genres WHERE movie_id = $1", id)
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

	// Update casts
	if update.Casts != nil {
		_, err := dbTx.Exec(ctx, "DELETE FROM movies_cast WHERE movie_id = $1", id)
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

	// insert/update schedules
	scheduleIDMap := map[string]int{}
	keepKeys := map[string]struct{}{}

	if update.Schedules != nil {
		for _, s := range *update.Schedules {
			key := s.Date + "|" + s.Time
			keepKeys[key] = struct{}{}

			var scheduleID int
			err := dbTx.QueryRow(ctx,
				"SELECT id FROM schedules WHERE movie_id=$1 AND date=$2 AND time=$3",
				id, s.Date, s.Time,
			).Scan(&scheduleID)

			if err != nil {
				if err == pgx.ErrNoRows {
					err = dbTx.QueryRow(ctx,
						"INSERT INTO schedules (movie_id, date, time) VALUES ($1, $2, $3) RETURNING id",
						id, s.Date, s.Time,
					).Scan(&scheduleID)
					if err != nil {
						return err
					}
				} else {
					return err
				}
			}
			scheduleIDMap[key] = scheduleID
		}
	}

	// remove schedules base on payload data
	rows, err := dbTx.Query(ctx, "SELECT id, date::text, time::text FROM schedules WHERE movie_id=$1", id)
	if err != nil {
		return err
	}

	var dbSchedules []models.ScheduleDB

	for rows.Next() {
		var s models.ScheduleDB
		if err := rows.Scan(&s.ID, &s.Date, &s.Time); err != nil {
			rows.Close()
			return err
		}
		dbSchedules = append(dbSchedules, s)
	}
	rows.Close()

	for _, s := range dbSchedules {
		key := s.Date + "|" + s.Time
		if _, ok := keepKeys[key]; !ok {
			if _, err := dbTx.Exec(ctx, "DELETE FROM cinemas_schedules WHERE schedules_id=$1", s.ID); err != nil {
				return err
			}
			if _, err := dbTx.Exec(ctx, "DELETE FROM schedules WHERE id=$1", s.ID); err != nil {
				return err
			}
		}
	}

	// Insert cinema schedules
	if update.CinemaSchedules != nil {
		for _, cs := range *update.CinemaSchedules {
			timeKey := cs.Time
			if len(timeKey) > 5 {
				timeKey = timeKey[:5]
			}
			key := cs.Date + "|" + timeKey
			scheduleID, ok := scheduleIDMap[key]
			if !ok {
				log.Printf("[WARN] Schedule not found for CinemaSchedule %+v\n", cs)
				continue
			}

			_, err := dbTx.Exec(ctx,
				`INSERT INTO cinemas_schedules (cinemas_id, schedules_id, locations_id)
                 VALUES ($1, $2, $3)
                 ON CONFLICT DO NOTHING`,
				cs.CinemaID, scheduleID, cs.LocationID,
			)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Cinema schedule inserted: %+v with scheduleID=%d\n", cs, scheduleID)
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
