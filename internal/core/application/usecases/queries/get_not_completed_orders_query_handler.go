package queries

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	"gorm.io/gorm"
)

type IGetNotCompletedOrdersHandler interface {
	Handle(GetNotCompletedOrders) (GetNotCompletedOrdersResponse, error)
}

type GetNotCompletedOrdersHandler struct {
	db *gorm.DB
}

func NewGetNotCompletedOrdersHandler(db *gorm.DB) (*GetNotCompletedOrdersHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &GetNotCompletedOrdersHandler{db: db}, nil
}

func (q *GetNotCompletedOrdersHandler) Handle(query GetNotCompletedOrders) (GetNotCompletedOrdersResponse, error) {
	if query.IsEmpty() {
		return GetNotCompletedOrdersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var orders []OrderResponse
	result := q.db.Raw(
		"SELECT id, courier_id, location_x, location_y, status FROM public.orders where status!=?",
		order.StatusCompleted,
	).Scan(&orders)

	if result.Error != nil {
		return GetNotCompletedOrdersResponse{}, result.Error
	}

	return GetNotCompletedOrdersResponse{Orders: orders}, nil
}
