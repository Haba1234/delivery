package problems

import (
	"errors"
	"net/http"
)

var ErrProblemBadRequest = errors.New("bad request")

type BadRequestError struct {
	ProblemDetailsError
}

func NewBadRequestError(detail string) *BadRequestError {
	return &BadRequestError{
		ProblemDetailsError: ProblemDetailsError{
			Type:   "bad-request",
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: detail,
		},
	}
}

func (e *BadRequestError) Error() string {
	return e.ProblemDetailsError.Error()
}

func (*BadRequestError) Unwrap() error {
	return ErrProblemBadRequest
}
