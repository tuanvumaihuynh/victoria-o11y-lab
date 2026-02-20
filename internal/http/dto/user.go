package dto

import "time"

type CreateUserPayload struct {
	Name  string `json:"name" minLength:"1" example:"John Doe"`
	Email string `json:"email" format:"email" example:"john.doe@example.com"`
	//nolint:gosec
	Password string `json:"password" minLength:"8" example:"password123"`
}

type CreateUserRequest struct {
	Body CreateUserPayload
}

type CreateUserResponseBody struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-01-01T00:00:00Z"`
}

type CreateUserResponse struct {
	Body CreateUserResponseBody
}

type GetUserByIDRequest struct {
	ID string `path:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type GetUserByIDResponseBody CreateUserResponseBody

type GetUserByIDResponse struct {
	Body GetUserByIDResponseBody
}
