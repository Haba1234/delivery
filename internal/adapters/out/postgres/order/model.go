package order

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
)

type ModelOrder struct {
	ID        uuid.UUID     `gorm:"type:uuid;primaryKey"`
	CourierID *uuid.UUID    `gorm:"type:uuid;index"`
	Location  ModelLocation `gorm:"embedded;embeddedPrefix:location_"`
	Status    order.Status  `gorm:"type:varchar(20)"`
}

type ModelLocation struct {
	X int
	Y int
}

func (ModelOrder) TableName() string {
	return "orders"
}
