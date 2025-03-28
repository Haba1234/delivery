package courier

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelCourier struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Transport ModelTransport `gorm:"foreignKey:CourierID;references:ID;constraint:OnDelete:CASCADE;"`
	Location  ModelLocation  `gorm:"embedded;embeddedPrefix:location_"`
	Status    courier.Status `gorm:"type:varchar(20)"`
}

type ModelLocation struct {
	X int
	Y int
}

type ModelTransport struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Speed     int
	CourierID uuid.UUID `gorm:"type:uuid;index"`
}

func (ModelCourier) TableName() string {
	return "couriers"
}

func (ModelTransport) TableName() string {
	return "transports"
}
