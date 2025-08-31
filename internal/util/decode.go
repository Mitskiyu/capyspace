package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator interface {
	Valid(ctx context.Context) map[string]string
}

func Decode[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("failed to decode json: %w", err)
	}

	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("failed to validate %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}
