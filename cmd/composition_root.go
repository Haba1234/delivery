package cmd

import (
	"context"
	"log"

	"github.com/Haba1234/delivery/internal/adapters/in/jobs"
	kafkain "github.com/Haba1234/delivery/internal/adapters/in/kafka"
	"github.com/Haba1234/delivery/internal/adapters/out/grpc/geo"
	kafkaout "github.com/Haba1234/delivery/internal/adapters/out/kafka"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/application/eventhandlers"
	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/core/application/usecases/queries"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/domain/services"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/uow"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/robfig/cron/v3"
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

type Jobs struct {
	AssignOrders cron.Job
	MoveCouriers cron.Job
}

type Clients struct {
	GeoClient ports.IGeoClient
}

type Consumers struct {
	OrdersCreateConsumer kafkain.IOrdersCreateConsumer
}

type Producers struct {
	OrderConfirmedProducer ports.IOrderProducer
}

type CompositionRoot struct {
	DomainServices  DomainServices
	Repositories    Repositories
	CommandHandlers CommandHandlers
	QueryHandlers   QueryHandlers
	Jobs            Jobs
	Clients         Clients
	Consumers       Consumers
	Producers       Producers

	closeFns []func() error
}

func NewCompositionRoot(_ context.Context, db *gorm.DB, configs Configs) CompositionRoot {
	// Domain Services
	orderDispatcher := services.NewDispatchService()

	// Repositories
	unitOfWork, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	orderRepository, err := orderrepo.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	courierRepository, err := courier.NewRepository(db)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Grpc Clients
	geoClient, err := geo.NewClient(configs.GeoClientURL)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Kafka Producers
	orderKafkaProducer, err := kafkaout.NewOrderProducer(
		configs.KafkaHost, configs.KafkaOrdersStatusChangedTopic,
	)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Command Handlers
	createOrderHandler, err := commands.NewCreateOrderHandler(orderRepository, geoClient)
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

	// Domain Event Handlers
	orderDomainEventHandler, err := eventhandlers.NewOrderCompletedDomainEventHandler(orderKafkaProducer)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Mediatr Subscribes
	err = mediatr.RegisterNotificationHandlers[order.CompletedDomainEvent](orderDomainEventHandler)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Kafka Consumers
	ordersCreateConsumer, err := kafkain.NewOrdersCreateConsumer(
		configs.KafkaHost, configs.ConsumerGroup,
		configs.KafkaOrdersCreateTopic, createOrderHandler,
	)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	// Jobs
	assignOrdersJob, err := jobs.NewAssignOrders(assignOrdersHandler)
	if err != nil {
		log.Fatalf("run application error: %s", err)
	}

	moveCouriersJob, err := jobs.NewMoveCouriers(moveCouriersHandler)
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
		Jobs: Jobs{
			assignOrdersJob,
			moveCouriersJob,
		},
		Consumers: Consumers{
			OrdersCreateConsumer: ordersCreateConsumer,
		},
		Producers: Producers{
			OrderConfirmedProducer: orderKafkaProducer,
		},
	}

	// Close
	compositionRoot.closeFns = append(compositionRoot.closeFns, geoClient.Close)
	compositionRoot.closeFns = append(compositionRoot.closeFns, ordersCreateConsumer.Close)
	compositionRoot.closeFns = append(compositionRoot.closeFns, orderKafkaProducer.Close)

	return compositionRoot
}

func (cr *CompositionRoot) Close() {
	for _, fn := range cr.closeFns {
		if err := fn(); err != nil {
			log.Printf("ошибка при закрытии зависимости: %v", err)
		}
	}
}
