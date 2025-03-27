package queries

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	gorm.Model
	ID       uuid.UUID
	Location ModelLocation `gorm:"embedded;embeddedPrefix:location_"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
