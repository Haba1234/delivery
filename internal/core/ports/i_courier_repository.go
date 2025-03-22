package ports

import (
	"context"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=./../../mocks/courier_repository_mock.go -package=mocks . ICourierRepository
type ICourierRepository interface {
	Add(ctx context.Context, aggregate *courier.Courier) error
	Update(ctx context.Context, aggregate *courier.Courier) error
	Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error)
	GetAllInFreeStatus(ctx context.Context) ([]*courier.Courier, error)
}
