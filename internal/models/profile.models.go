package models

type Profile struct {
	ID             int     `json:"id"`
	Email          string  `json:"email"`
	Password       string  `json:"-"`
	Role           string  `json:"role"`
	VirtualAccount string  `json:"virtual_account"`
	FirstName      *string `json:"first_name"`
	LastName       *string `json:"last_name"`
	PhoneNumber    *string `json:"phone_number"`
	Points         *int    `json:"points"`
	ImagePath      *string `json:"image_path"`
}

type UserUpdate struct {
	Email       *string `form:"email" json:"email"`
	OldPassword *string `form:"old_password" json:"old_password"`
	NewPassword *string `form:"new_password" json:"new_password"`
	Password    *string `json:"-"`
}

type ProfileUpdate struct {
	FirstName   *string `form:"first_name" json:"first_name"`
	LastName    *string `form:"last_name" json:"last_name"`
	PhoneNumber *string `form:"phone_number" json:"phone_number"`
	ImagePath   *string `form:"image_path" json:"image_path"`
	Points      *int    `form:"points" json:"points"`
}

type UserUpdateRequest struct {
	User    UserUpdate    `json:"user"`
	Profile ProfileUpdate `json:"profile"`
}
