package users_service

import (
	"context"
	"fmt"

	"github.com/MlkProduction/golang-project/internal/core/domain"
	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
)

func (s *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf("limit don't can between zero value: %w", core_error.ErrInvalidArgument)
	}

	if offset != nil && *limit < 0 {
		return nil, fmt.Errorf("offset don't can between zero value: %w", core_error.ErrInvalidArgument)
	}

	users, err := s.usersRepository.GetUsers(ctx, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil

}
