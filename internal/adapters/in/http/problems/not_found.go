package problems

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type NotFoundError struct {
	ProblemDetailsError
}

func NewNotFoundError(detail string) *NotFoundError {
	return &NotFoundError{
		ProblemDetailsError: ProblemDetailsError{
			Type:   "not-found",
			Title:  "Resource Not Found",
			Status: http.StatusNotFound,
			Detail: detail,
		},
	}
}

func (e *NotFoundError) Error() string {
	return e.ProblemDetailsError.Error()
}

func (*NotFoundError) Unwrap() error {
	return ErrNotFound
}
