package queries

import "github.com/google/uuid"

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID       uuid.UUID
	Location LocationResponse
}
