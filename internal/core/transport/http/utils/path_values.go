package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
)

/*
GET /users/{id}

Надо вычленить id из пути запроса
*/

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)

	if pathValue == "" {
		return 0, fmt.Errorf("no key: '%s' in path values: %w", key, core_error.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("Не получилось преобразовать: %s из ключа: %s в число: %v: %w", pathValue, key, err, core_error.ErrInvalidArgument)
	}

	return val, nil

}
