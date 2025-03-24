package cmd

import (
	"context"
	"log"

	"github.com/Haba1234/delivery/internal/adapters/out/postgres"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/core/application/usecases/queries"
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

type CommandHandlers struct {
	AssignOrderHandler  commands.IAssignOrderHandler
	CreateOrderHandler  commands.ICreateOrderHandler
	MoveCouriersHandler commands.IMoveCouriersHandler
}

type QueryHandlers struct {
	GetAllCouriersHandler        queries.IGetAllCouriersHandler
	GetNotCompletedOrdersHandler queries.IGetNotCompletedOrdersHandler
}

type CompositionRoot struct {
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
}

func NewCompositionRoot(_ context.Context, db *gorm.DB) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewDispatchService()

	// Repositories
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	orderRepository, err := order.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	courierRepository, err := courier.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Command Handlers
	createOrderHandler, err := commands.NewCreateOrderHandler(orderRepository)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	assignOrdersHandler, err := commands.NewAssignOrderHandler(
		unitOfWork, orderRepository, courierRepository, orderDispatcher,
	)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	moveCouriersHandler, err := commands.NewMoveCouriersHandler(
		unitOfWork, orderRepository, courierRepository,
	)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Query Handlers
	getAllCouriersHandler, err := queries.NewGetAllCouriersHandler(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	getNotCompletedOrdersHandler, err := queries.NewGetNotCompletedOrdersHandler(db)
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
		CommandHandlers: CommandHandlers{
			AssignOrderHandler:  assignOrdersHandler,
			CreateOrderHandler:  createOrderHandler,
			MoveCouriersHandler: moveCouriersHandler,
		},
		QueryHandlers: QueryHandlers{
			GetAllCouriersHandler:        getAllCouriersHandler,
			GetNotCompletedOrdersHandler: getNotCompletedOrdersHandler,
		},
	}

	return compositionRoot
}
