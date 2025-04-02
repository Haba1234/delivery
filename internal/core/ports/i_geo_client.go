package ports

import (
	"context"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
)

//go:generate mockgen -destination=./../../mocks/geo_client_mock.go -package=mocks . IGeoClient
type IGeoClient interface {
	GetGeolocation(ctx context.Context, street string) (kernel.Location, error)
}
