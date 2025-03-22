package commands

import (
	"context"
	"errors"

	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"
)

var ErrOrderAlreadyExists = errors.New("order already exists")

type ICreateOrderHandler interface {
	Handle(context.Context, CreateOrder) error
}

var _ ICreateOrderHandler = &CreateOrderHandler{}

type CreateOrderHandler struct {
	orderRepository ports.IOrderRepository
}

func NewCreateOrderHandler(orderRepository ports.IOrderRepository) (*CreateOrderHandler, error) {
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}

	return &CreateOrderHandler{
		orderRepository: orderRepository,
	}, nil
}

func (ch *CreateOrderHandler) Handle(ctx context.Context, command CreateOrder) error {
	if command.IsEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Проверяем нет ли уже такого заказа
	orderAggregate, err := ch.orderRepository.Get(ctx, command.OrderID())
	if err != nil && !errors.Is(err, orderrepo.ErrOrderNotFound) {
		return err
	}

	if orderAggregate != nil {
		return ErrOrderAlreadyExists
	}

	// Получили геопозицию
	location, err := kernel.CreateRandomLocation()
	if err != nil {
		return err
	}

	// Создали заказ
	orderAggregate, err = order.New(command.OrderID(), location)
	if err != nil {
		return err
	}

	// Сохранили
	return ch.orderRepository.Add(ctx, orderAggregate)
}
