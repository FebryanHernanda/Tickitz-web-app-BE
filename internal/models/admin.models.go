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
	Title        string     `form:"title" json:"title" example:"Negeri 5 Menara"`
	PosterPath   string     `form:"poster_path" json:"poster_path" example:"/path/poster.jpg"`
	BackdropPath string     `form:"backdrop_path" json:"backdrop_path" example:"/path/backdrop.jpg"`
	Synopsis     string     `form:"synopsis" json:"synopsis" example:"Negeri 5 Menara merupakan film yang..."`
	ReleaseDate  string     `form:"release_date" json:"release_date" example:"2025-09-01"`
	Rating       float32    `form:"rating" json:"rating" example:"7.5"`
	AgeRating    string     `form:"age_rating" json:"age_rating" example:"R"`
	Duration     int        `form:"duration" json:"duration" example:"120"`
	DirectorID   int        `form:"director_id" json:"director_id" example:"1"`
	Genres       []int      `json:"genres" example:"1"`
	Casts        []int      `json:"casts" example:"2"`
	Schedules    []Schedule `json:"schedules"`
}

type EditMovies struct {
	Title        *string     `form:"title" json:"title,omitempty" example:"Negeri 5 Menara"`
	PosterPath   *string     `form:"poster_path" json:"poster_path,omitempty" example:"/path/poster.jpg"`
	BackdropPath *string     `form:"backdrop_path" json:"backdrop_path,omitempty" example:"/path/backdrop.jpg"`
	Synopsis     *string     `form:"synopsis" json:"synopsis,omitempty" example:"Negeri 5 Menara merupakan film yang..."`
	ReleaseDate  *string     `form:"release_date" json:"release_date,omitempty" example:"2025-09-01"`
	Rating       *float32    `form:"rating" json:"rating,omitempty" example:"7.5"`
	AgeRating    *string     `form:"age_rating" json:"age_rating,omitempty" example:"R"`
	Duration     *int        `form:"duration" json:"duration,omitempty" example:"120"`
	DirectorID   *int        `form:"director_id" json:"director_id,omitempty" example:"1"`
	Genres       *[]int      `json:"genres,omitempty" example:"1"`
	Casts        *[]int      `json:"casts,omitempty" example:"2"`
	Schedules    *[]Schedule `json:"schedules,omitempty"`
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
