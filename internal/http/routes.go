package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (s *Service) RegisterRoutes(api huma.API) {
	group := huma.NewGroup(api, "/api/v1")

	registerHandler(group, http.MethodPost, "/users", s.CreateUser, CreateUserDocs())
	registerHandler(group, http.MethodGet, "/users/{id}", s.GetUserByID, GetUserByIDDocs())
}

func registerHandler[I any, O any](
	humaAPI huma.API,
	method string,
	path string,
	handler func(ctx context.Context, req *I) (*O, error),
	op huma.Operation,
) {
	op.Method = method
	op.Path = path
	huma.Register(humaAPI, op, func(ctx context.Context, req *I) (*O, error) {
		output, err := handler(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("handler error: %w", err)
		}

		return output, nil
	})
}
