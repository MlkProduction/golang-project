package users_transport_http

import (
	"context"
	"net/http"

	"github.com/MlkProduction/golang-project/internal/core/domain"
	core_http_server "github.com/MlkProduction/golang-project/internal/core/transport/http/server"
)

// UsersHTTPHandler — HTTP-обработчик для работы с пользователями.
// Хранит зависимость от сервиса, который содержит бизнес-логику работы с пользователями.
type UsersHTTPHandler struct {
	// usersService — сервис пользователей.
	// Через него HTTP-слой будет обращаться к бизнес-логике.
	usersService UsersService
}

// UsersService — интерфейс сервиса пользователей.
// Определяет методы, которые нужны HTTP-обработчику для работы с пользователями.
type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error)
}

// newUserHTTPHandler — конструктор UsersHTTPHandler.
// Принимает сервис пользователей и сохраняет его в обработчике.
func NewUserHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
