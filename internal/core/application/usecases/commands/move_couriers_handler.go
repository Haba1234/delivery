package commands

import (
	"context"
	"errors"

	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/internal/pkg/uow"
)

type IMoveCouriersHandler interface {
	Handle(context.Context, MoveCouriers) error
}

var _ IMoveCouriersHandler = &MoveCouriersHandler{}

type MoveCouriersHandler struct {
	unitOfWork        uow.IUnitOfWork
	orderRepository   ports.IOrderRepository
	courierRepository ports.ICourierRepository
}

func NewMoveCouriersHandler(
	unitOfWork uow.IUnitOfWork,
	orderRepository ports.IOrderRepository,
	courierRepository ports.ICourierRepository,
) (*MoveCouriersHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}
	if orderRepository == nil {
		return nil, errs.NewValueIsRequiredError("orderRepository")
	}
	if courierRepository == nil {
		return nil, errs.NewValueIsRequiredError("courierRepository")
	}

	return &MoveCouriersHandler{
		unitOfWork:        unitOfWork,
		orderRepository:   orderRepository,
		courierRepository: courierRepository,
	}, nil
}

func (ch *MoveCouriersHandler) Handle(ctx context.Context, command MoveCouriers) error {
	if command.IsEmpty() {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Читаем заказы в статусе "Назначен"
	assignedOrders, err := ch.orderRepository.GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, orderrepo.ErrAssignedOrdersNotFound) {
			return nil // Нет заказов в статусе "Назначен", ничего не делаем
		}

		return err
	}

	// Начинаем транзакцию
	ctx = ch.unitOfWork.Begin(ctx)
	for _, assignedOrder := range assignedOrders {
		if err := ch.processAssignedOrder(ctx, assignedOrder); err != nil {
			ch.unitOfWork.Rollback(ctx)
			return err
		}
	}

	// Завершаем транзакцию
	return ch.unitOfWork.Commit(ctx)
}

// processAssignedOrder обрабатывает назначенный заказ и перемещает курьера на его местонахождение.
func (ch *MoveCouriersHandler) processAssignedOrder(ctx context.Context, o *order.Order) error {
	// Получаем курьера из заказа
	courier, err := ch.courierRepository.Get(ctx, *o.CourierID())
	if err != nil {
		return err
	}

	// Перемещаем курьера к заказу
	if err := courier.Move(o.Location()); err != nil {
		return err
	}

	// Если местонахождение совпадает, завершаем заказ и освобождаем курьера
	if courier.Location().Equals(o.Location()) {
		if err := o.Complete(); err != nil {
			return err
		}
		if err := courier.SetFree(); err != nil {
			return err
		}
	}

	// Сохраняем изменения
	if err := ch.orderRepository.Update(ctx, o); err != nil {
		return err
	}

	if err := ch.courierRepository.Update(ctx, courier); err != nil {
		return err
	}

	return nil
}
