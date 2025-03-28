package queries

import (
	"github.com/google/uuid"
)

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID        uuid.UUID `gorm:"column:id"`
	LocationX int       `gorm:"column:location_x"`
	LocationY int       `gorm:"column:location_y"`
}
