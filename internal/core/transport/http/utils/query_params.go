package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("Не получилось преобразовать: %s из ключа: %s в число: %v: %w", param, key, err, core_error.ErrInvalidArgument)
	}

	return &val, nil

}
