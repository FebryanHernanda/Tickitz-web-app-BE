package models

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

type UserUpdate struct {
	Email    *string `form:"email" json:"email"`
	Password *string `form:"password" json:"password"`
}

type ProfileUpdate struct {
	FirstName   *string `form:"first_name" json:"first_name"`
	LastName    *string `form:"last_name" json:"last_name"`
	PhoneNumber *string `form:"phone_number" json:"phone_number"`
	ImagePath   *string `json:"image_path"`
	Points      *int    `json:"points"`
}

type UserUpdateRequest struct {
	User    UserUpdate    `json:"user"`
	Profile ProfileUpdate `json:"profile"`
}
