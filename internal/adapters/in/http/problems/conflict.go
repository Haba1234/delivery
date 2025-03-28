package problems

import (
	"errors"
	"net/http"
)

var ErrProblemConflict = errors.New("conflict")

type ConflictError struct {
	ProblemDetailsError
}

func NewConflictError(problemType, detail string) *ConflictError {
	return &ConflictError{
		ProblemDetailsError: ProblemDetailsError{
			Type:   problemType,
			Title:  "Conflict",
			Status: http.StatusConflict,
			Detail: detail,
		},
	}
}

func (e *ConflictError) Error() string {
	return e.ProblemDetailsError.Error()
}

func (*ConflictError) Unwrap() error {
	return ErrProblemConflict
}
