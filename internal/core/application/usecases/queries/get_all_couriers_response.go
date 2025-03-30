package queries

import (
	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID        uuid.UUID `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	LocationX int       `gorm:"column:location_x"`
	LocationY int       `gorm:"column:location_y"`
}
