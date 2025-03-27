package http

import (
	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/core/application/usecases/queries"
	"github.com/Haba1234/delivery/internal/pkg/errs"
)

type Server struct {
	createOrderCommandHandler         commands.ICreateOrderHandler
	getAllCouriersQueryHandler        queries.IGetAllCouriersHandler
	getNotCompletedOrdersQueryHandler queries.IGetNotCompletedOrdersHandler
}

func NewServer(
	createOrderCommandHandler commands.ICreateOrderHandler,

	getAllCouriersQueryHandler queries.IGetAllCouriersHandler,
	getNotCompletedOrdersQueryHandler queries.IGetNotCompletedOrdersHandler,
) (*Server, error) {
	if createOrderCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("createOrderCommandHandler")
	}

	if getAllCouriersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getAllCouriersQueryHandler")
	}

	if getNotCompletedOrdersQueryHandler == nil {
		return nil, errs.NewValueIsRequiredError("getNotCompletedOrdersQueryHandler")
	}

	return &Server{
		createOrderCommandHandler:         createOrderCommandHandler,
		getAllCouriersQueryHandler:        getAllCouriersQueryHandler,
		getNotCompletedOrdersQueryHandler: getNotCompletedOrdersQueryHandler,
	}, nil
}
