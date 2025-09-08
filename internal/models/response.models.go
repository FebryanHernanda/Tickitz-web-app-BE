package models

type ErrorResponse struct {
	Success bool   `example:"false"`
	Error   string `example:"error message"`
}

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}
