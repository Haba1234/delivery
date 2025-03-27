package queries

import (
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IGetAllCouriersHandler interface {
	Handle(GetAllCouriers) (GetAllCouriersResponse, error)
}

type GetAllCouriersHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersHandler(db *gorm.DB) (*GetAllCouriersHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &GetAllCouriersHandler{db: db}, nil
}

func (q *GetAllCouriersHandler) Handle(query GetAllCouriers) (GetAllCouriersResponse, error) {
	if query.IsEmpty() {
		return GetAllCouriersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var couriers []CourierResponse

	result := q.db.
		Preload(clause.Associations).
		Find(&couriers)

	if result.Error != nil {
		return GetAllCouriersResponse{}, result.Error
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}
