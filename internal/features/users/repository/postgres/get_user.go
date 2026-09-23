package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MlkProduction/golang-project/internal/core/domain"
	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())

	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number 
	FROM todoapp.users
	WHERE id=$1
	`
	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d': %w", id, core_error.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("get userModel: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)

	return userDomain, nil

}
