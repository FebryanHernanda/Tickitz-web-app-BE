package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	Password       string    `json:"password"`
	VirtualAccount string    `json:"virtual_account"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `json:"password" binding:"required,min=8" example:"User@testing123"`
	Role     string `json:"role,omitempty" example:"user"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email" example:"user@mail.com"`
	Password string `json:"password" binding:"required" example:"your_password"`
}

type Profile struct {
	ID             int     `json:"id"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	VirtualAccount string  `json:"virtual_account"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	PhoneNumber    *string `json:"phone_number"`
	Points         *int    `json:"points"`
	ImagePath      *string `json:"image_path"`
}
