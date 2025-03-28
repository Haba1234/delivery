package problems

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ProblemDetailsError RFC 7807.
type ProblemDetailsError struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (p *ProblemDetailsError) Error() string {
	return fmt.Sprintf("%d: %s - %s", p.Status, p.Title, p.Detail)
}

func (p *ProblemDetailsError) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}
