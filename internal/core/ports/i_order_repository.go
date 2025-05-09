package ports

import (
	"context"

	"github.com/Haba1234/delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=./../../mocks/order_repository_mock.go -package=mocks . IOrderRepository
type IOrderRepository interface {
	Add(ctx context.Context, aggregate *order.Order) error
	Update(ctx context.Context, aggregate *order.Order) error
	Get(ctx context.Context, ID uuid.UUID) (*order.Order, error)
	GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error)
	GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error)
}
