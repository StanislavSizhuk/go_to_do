package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/StanislavSizhuk/go_to_do/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {

	pathValues := r.PathValue(key)
	if pathValues == "" {
		return 0, fmt.Errorf("no key '%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValues)
	if err != nil {
		return 0, fmt.Errorf("path value '%s' by key  '%s' is not a valid integer: %v: %w", pathValues, key, err, core_errors.ErrInvalidArgument)
	}
	return val, nil
}
