package repositories

import (
	"context"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MoviesRepository struct {
	DB *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{
		DB: db,
	}
}

func (mr *MoviesRepository) GetGenreMovies(ctx context.Context) ([]models.MoviesGenres, error) {
	query := `SELECT id,name FROM genres`

	rows, err := mr.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.MoviesGenres
	for rows.Next() {
		var mv models.MoviesGenres
		err := rows.Scan(
			&mv.ID,
			&mv.Name,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, mv)
	}
	return genres, nil
}

func (mr *MoviesRepository) GetCastsMovies(ctx context.Context) ([]models.MoviesCast, error) {
	query := `SELECT id,name FROM casts`

	rows, err := mr.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var casts []models.MoviesCast
	for rows.Next() {
		var mv models.MoviesCast
		err := rows.Scan(
			&mv.ID,
			&mv.Name,
		)
		if err != nil {
			return nil, err
		}
		casts = append(casts, mv)
	}
	return casts, nil
}
func (mr *MoviesRepository) GetDirectorsMovies(ctx context.Context) ([]models.MoviesDirectors, error) {
	query := `SELECT id,name FROM directors`

	rows, err := mr.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directors []models.MoviesDirectors
	for rows.Next() {
		var mv models.MoviesDirectors
		err := rows.Scan(
			&mv.ID,
			&mv.Name,
		)
		if err != nil {
			return nil, err
		}
		directors = append(directors, mv)
	}
	return directors, nil
}

func (mr *MoviesRepository) GetUpcomingMovies(ctx context.Context) ([]models.Movie, error) {
	query := `
	SELECT
		m.id,
		m.title,
		m.poster_path,
		m.backdrop_path,
		m.release_date,
    	ARRAY_AGG (g.name)
	FROM
		movies m
		JOIN movies_genres mg ON m.id = mg.movie_id
		JOIN genres g ON mg.genre_id = g.id
	WHERE
    	release_date > CURRENT_DATE
	GROUP BY
    	m.id;
	`

	rows, err := mr.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var mv models.Movie
		err := rows.Scan(
			&mv.ID,
			&mv.Title,
			&mv.PosterPath,
			&mv.BackdropPath,
			&mv.ReleaseDate,
			&mv.Genres,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, mv)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetPopularMovies(ctx context.Context) ([]models.Movie, error) {
	query := `
		SELECT
			m.id,
			m.title,
			m.poster_path,
			m.backdrop_path,
			m.release_date,
			ARRAY_AGG (g.name)
		FROM
			movies m
			JOIN movies_genres mg ON m.id = mg.movie_id
			JOIN genres g ON mg.genre_id = g.id
		WHERE
			rating > 8
		GROUP BY
			m.id
		LIMIT
			10;
	`
	rows, err := mr.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var mv models.Movie
		err := rows.Scan(
			&mv.ID,
			&mv.Title,
			&mv.PosterPath,
			&mv.BackdropPath,
			&mv.ReleaseDate,
			&mv.Genres,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, mv)
	}
	return movies, nil
}

func (mr *MoviesRepository) GetMoviesByFilter(ctx context.Context, search, genre string, page, limit, offset int) ([]models.Movie, int, error) {
	var totalCount int
	countQuery := `
        SELECT COUNT(DISTINCT m.id)
        FROM movies m
        JOIN movies_genres mg ON m.id = mg.movie_id
        JOIN genres g ON mg.genre_id = g.id
        WHERE 
            ($1 = '' OR m.title ILIKE '%' || $1 || '%')
            AND ($2 = '' OR g.name ILIKE '%' || $2 || '%');
    `
	err := mr.DB.QueryRow(ctx, countQuery, search, genre).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `
        SELECT
            m.id,
            m.title,
            m.poster_path,
            m.backdrop_path,
            m.release_date,
            ARRAY_AGG(g.name) AS genres
        FROM movies m
        JOIN movies_genres mg ON m.id = mg.movie_id
        JOIN genres g ON mg.genre_id = g.id
        WHERE 
            ($1 = '' OR m.title ILIKE '%' || $1 || '%')
            AND ($2 = '' OR g.name ILIKE '%' || $2 || '%')
        GROUP BY m.id
        ORDER BY m.id ASC
        LIMIT $3 OFFSET $4;
    `
	values := []any{search, genre, limit, offset}

	rows, err := mr.DB.Query(ctx, query, values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var mv models.Movie
		err := rows.Scan(
			&mv.ID,
			&mv.Title,
			&mv.PosterPath,
			&mv.BackdropPath,
			&mv.ReleaseDate,
			&mv.Genres,
		)
		if err != nil {
			return nil, 0, err
		}
		movies = append(movies, mv)
	}
	return movies, totalCount, nil
}

func (mr *MoviesRepository) GetDetailMovies(ctx context.Context, movieID int64) (*models.MovieDetails, error) {
	query := `
	SELECT
		m.id,
		m.title,
		m.poster_path,
		m.backdrop_path,
		m.release_date,
		m.rating,
		m.duration,
		m.synopsis,
		d.name AS director,
		ARRAY_AGG (DISTINCT c.name) AS casts,
		ARRAY_AGG (DISTINCT g.name) AS genres
	FROM
		movies m
		LEFT JOIN directors d ON m.director_id = d.id
		LEFT JOIN movies_cast mc ON m.id = mc.movie_id
		LEFT JOIN casts c ON mc.cast_id = c.id
		LEFT JOIN movies_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
	WHERE
    	m.id = $1
	GROUP BY
		m.id,
		d.name
	`

	row := mr.DB.QueryRow(ctx, query, movieID)

	var mv models.MovieDetails
	err := row.Scan(
		&mv.ID,
		&mv.Title,
		&mv.PosterPath,
		&mv.BackdropPath,
		&mv.ReleaseDate,
		&mv.Rating,
		&mv.Duration,
		&mv.Synopsis,
		&mv.Director,
		&mv.Casts,
		&mv.Genres,
	)
	if err != nil {
		return nil, err
	}

	return &mv, err
}

func (mr *MoviesRepository) GetSchedulesMovies(ctx context.Context) ([]models.MovieSchedules, error) {
	query := `
        SELECT 
            s.id,
            s.date,
            s.time,
            m.title,
            c.name,
			c.prices AS ticket_price,
            l.name
        FROM schedules s
        JOIN movies m ON m.id = s.movie_id
        JOIN cinemas_schedules cs ON cs.schedules_id = s.id
        JOIN cinemas c ON c.id = cs.cinemas_id
        JOIN locations l ON l.id = cs.locations_id
    `

	rows, err := mr.DB.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.MovieSchedules
	for rows.Next() {
		var sch models.MovieSchedules
		err := rows.Scan(
			&sch.ID,
			&sch.Date,
			&sch.Time,
			&sch.MovieTitle,
			&sch.CinemaName,
			&sch.TicketPrices,
			&sch.LocationName,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, sch)
	}
	return schedules, nil
}
