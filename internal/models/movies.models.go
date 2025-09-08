package models

import "time"

type Movie struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
	ReleaseDate  time.Time `json:"release_date"`
	Genres       []string  `json:"genres"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MovieDetails struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
	ReleaseDate  time.Time `json:"release_date"`
	Rating       float64   `json:"rating"`
	Duration     int       `json:"duration"`
	Synopsis     string    `json:"synopsis"`
	Director     string    `json:"director"`
	Casts        []string  `json:"casts"`
	Genres       []string  `json:"genres"`
}

type MovieSchedules struct {
	ID           int       `json:"id"`
	Date         time.Time `json:"date"`
	Time         string    `json:"time"`
	MovieTitle   string    `json:"movie_title"`
	CinemaName   string    `json:"cinema_name"`
	LocationName string    `json:"location_name"`
}
