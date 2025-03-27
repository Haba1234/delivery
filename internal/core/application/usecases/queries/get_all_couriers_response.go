package queries

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	gorm.Model
	ID       uuid.UUID
	Name     string
	Location ModelLocation `gorm:"embedded;embeddedPrefix:location_"`
}

func (CourierResponse) TableName() string {
	return "couriers"
}
