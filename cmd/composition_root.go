package cmd

import (
	"context"
	"log"

	"github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/services"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/uow"

	"gorm.io/gorm"
)

type DomainServices struct {
	OrderDispatcher services.IOrderDispatcher
}

type Repositories struct {
	UnitOfWork        uow.IUnitOfWork
	OrderRepository   ports.IOrderRepository
	CourierRepository ports.ICourierRepository
}

type CompositionRoot struct {
	DomainServices DomainServices
	Repositories   Repositories
}

func NewCompositionRoot(_ context.Context, db *gorm.DB) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewDispatchService()

	orderRepository, err := order.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	courierRepository, err := courier.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	compositionRoot := CompositionRoot{
		DomainServices: DomainServices{
			OrderDispatcher: orderDispatcher,
		},

		Repositories: Repositories{
			OrderRepository:   orderRepository,
			CourierRepository: courierRepository,
		},
	}

	return compositionRoot
}
