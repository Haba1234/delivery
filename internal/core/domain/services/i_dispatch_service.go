package services

import (
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
)

type IOrderDispatcher interface {
	Dispatch(*order.Order, []*courier.Courier) (*courier.Courier, error)
}
