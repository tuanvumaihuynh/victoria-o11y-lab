package http

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/dto"
)

func CreateUserDocs() huma.Operation {
	return huma.Operation{
		OperationID:   "create-user",
		Summary:       "Create a new user",
		Description:   "Create a new user with the given name and email",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"users"},
	}
}

func (s *Service) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	now := time.Now()
	return &dto.CreateUserResponse{
		Body: dto.CreateUserResponseBody{
			ID:        "123e4567-e89b-12d3-a456-426614174000",
			Name:      req.Body.Name,
			Email:     req.Body.Email,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}

func GetUserByIDDocs() huma.Operation {
	return huma.Operation{
		OperationID:   "get-user-by-id",
		Summary:       "Get a user by ID",
		Description:   "Get a user by ID with the given ID",
		DefaultStatus: http.StatusOK,
		Tags:          []string{"users"},
	}
}

func (s *Service) GetUserByID(ctx context.Context, req *dto.GetUserByIDRequest) (*dto.GetUserByIDResponse, error) {
	now := time.Now()
	return &dto.GetUserByIDResponse{
		Body: dto.GetUserByIDResponseBody{
			ID:        req.ID,
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}, nil
}
