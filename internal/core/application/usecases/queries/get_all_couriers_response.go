package queries

import "github.com/google/uuid"

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID
	Name     string
	Location LocationResponse
}
