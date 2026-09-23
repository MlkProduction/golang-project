package users_postgres_repository

import (
	"context"
	"fmt"

	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, id int) error {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OptTimeout())

	defer cancel()

	query := `
	DELETE from todoapp.users
	WHERE id=$1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("User with id=`%d` : %w", id, core_error.ErrNotFound)
	}

	return nil

}
