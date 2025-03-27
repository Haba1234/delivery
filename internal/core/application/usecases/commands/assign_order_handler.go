package commands

import (
	"context"
	"errors"

	courierrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/services"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/internal/pkg/uow"
)

var (
	ErrNotAvailableOrders   = errors.New("not available orders")
	ErrNotAvailableCouriers = errors.New("not available couriers")
)

type IAssignOrderHandler interface {
	Handle(context.Context, AssignOrder) error
}

var _ IAssignOrderHandler = &AssignOrderHandler{}

type AssignOrderHandler struct {
	unitOfWork        uow.IUnitOfWork
	orderRepository   ports.IOrderRepository
	courierRepository ports.ICourierRepository
	orderDispatcher   services.IOrderDispatcher
}

func NewAssignOrderHandler(
	unitOfWork uow.IUnitOfWork,
	orderRepository ports.IOrderRepository,
	courierRepository ports.ICourierRepository,
	orderDispatcher services.IOrderDispatcher,
) (*AssignOrderHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}
	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}
	if orderDispatcher == nil {
		return nil, errs.NewValueIsRequiredError("orderDispatcher")
	}

	return &AssignOrderHandler{
		unitOfWork:        unitOfWork,
		orderRepository:   orderRepository,
		courierRepository: courierRepository,
		orderDispatcher:   orderDispatcher,
	}, nil
}

func (ch *AssignOrderHandler) Handle(ctx context.Context, command AssignOrder) error {
	if command.IsEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Ищем заказ в статусе "Создан"
	orderAggregate, err := ch.orderRepository.GetFirstInCreatedStatus(ctx)
	if err != nil {
		if errors.Is(err, orderrepo.ErrOrderCreatedNotFound) {
			return ErrNotAvailableOrders
		}
		return err
	}

	// Получаем всех свободных курьеров
	couriers, err := ch.courierRepository.GetAllInFreeStatus(ctx)
	if err != nil {
		if errors.Is(err, courierrepo.ErrNoFreeCouriers) {
			return ErrNotAvailableCouriers
		}
		return err
	}

	// Назначаем курьера
	assignedCourier, err := ch.orderDispatcher.Dispatch(orderAggregate, couriers)
	if err != nil {
		return err
	}

	if err = orderAggregate.Assign(assignedCourier); err != nil {
		return err
	}

	if err = assignedCourier.SetBusy(); err != nil {
		return err
	}

	// Сохраняем изменения
	ctx = ch.unitOfWork.Begin(ctx)

	if err = ch.orderRepository.Update(ctx, orderAggregate); err != nil {
		ch.unitOfWork.Rollback(ctx)
		return err
	}
	if err = ch.courierRepository.Update(ctx, assignedCourier); err != nil {
		ch.unitOfWork.Rollback(ctx)
		return err
	}

	return ch.unitOfWork.Commit(ctx)
}
