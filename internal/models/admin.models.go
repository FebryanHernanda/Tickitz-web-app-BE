package models

import "time"

type AdminMovies struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
	Synopsis     string    `json:"synopsis"`
	ReleaseDate  time.Time `json:"release_date"`
	Rating       float64   `json:"rating"`
	AgeRating    string    `json:"age_rating"`
	Duration     int       `json:"duration"`
	Director     string    `json:"director"`
	DatePlaying  time.Time `json:"date_playing"`
	LocationName *string   `json:"location_name,omitempty"`
	CinemaName   *string   `json:"cinema_name,omitempty"`
	TimePlaying  []string  `json:"time_playing"`
	Casts        []string  `json:"casts"`
	Genres       []string  `json:"genres"`
}

type AddMovies struct {
	Title        string     `json:"title" example:"Negeri 5 Menara"`
	PosterPath   string     `json:"poster_path" example:"/path/poster.jpg"`
	BackdropPath string     `json:"backdrop_path" example:"/path/backdrop.jpg"`
	Synopsis     string     `json:"synopsis" example:"Negeri 5 Menara merupakan film yang..."`
	ReleaseDate  string     `json:"release_date" example:"2025-09-01"`
	Rating       float32    `json:"rating" example:"7.5"`
	AgeRating    string     `json:"age_rating" example:"R"`
	Duration     int        `json:"duration" example:"120"`
	DirectorID   int        `json:"director_id" example:"1"`
	Genres       []int      `json:"genres" example:"1"`
	Casts        []int      `json:"casts" example:"2"`
	Schedules    []Schedule `json:"schedules"`
}

type Schedule struct {
	Date string `json:"date" example:"2025-09-10"`
	Time string `json:"time" example:"18:00"`
}

type CinemaScheduleLocation struct {
	CinemaID   int `json:"cinema_id"`
	ScheduleID int `json:"schedule_id"`
	LocationID int `json:"location_id"`
}

type GetSchedule struct {
	ID      int       `json:"id" example:"1"`
	Date    time.Time `json:"date" example:"2025-09-10"`
	Time    string    `json:"time" example:"18:00"`
	MovieID int       `json:"movie_id" example:"1"`
}
