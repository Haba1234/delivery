package commands

import (
	"testing"

	courierrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/core/domain/services"
	"github.com/Haba1234/delivery/internal/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAssignOrderHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockIOrderRepository(ctrl)
	mockCourierRepo := mocks.NewMockICourierRepository(ctrl)
	mockUnitOfWork := mocks.NewMockIUnitOfWork(ctrl)

	dispatcher := services.NewDispatchService()

	handler, err := NewAssignOrderHandler(mockUnitOfWork, mockOrderRepo, mockCourierRepo, dispatcher)
	require.NoError(t, err)

	t.Run(
		"empty command", func(t *testing.T) {
			cmd := AssignOrder{}
			err := handler.Handle(t.Context(), cmd)
			require.Error(t, err)
			assert.EqualError(t, err, "value is required: add address command")
		},
	)

	t.Run(
		"no available orders", func(t *testing.T) {
			mockOrderRepo.EXPECT().
				GetFirstInCreatedStatus(gomock.Any()).
				Return(nil, orderrepo.ErrOrderCreatedNotFound)

			cmd, err := NewAssignOrder()
			require.NoError(t, err)

			err = handler.Handle(t.Context(), cmd)
			assert.ErrorIs(t, err, ErrNotAvailableOrders)
		},
	)

	t.Run(
		"no available couriers", func(t *testing.T) {
			orderMock := &order.Order{}
			mockOrderRepo.EXPECT().
				GetFirstInCreatedStatus(gomock.Any()).
				Return(orderMock, nil)

			mockCourierRepo.EXPECT().
				GetAllInFreeStatus(gomock.Any()).
				Return(nil, courierrepo.ErrNoFreeCouriers)

			cmd, err := NewAssignOrder()
			require.NoError(t, err)

			err = handler.Handle(t.Context(), cmd)
			assert.ErrorIs(t, err, ErrNotAvailableCouriers)
		},
	)

	t.Run(
		"successful assignment", func(t *testing.T) {
			ctx := t.Context()

			orderLocation, err := kernel.CreateRandomLocation()
			require.NoError(t, err)
			orderMock, err := order.New(uuid.New(), orderLocation)
			require.NoError(t, err)

			courierLocation, err := kernel.CreateRandomLocation()
			require.NoError(t, err)

			courierMock, err := courier.New("Jon", "bike", 1, courierLocation)

			couriersMock := []*courier.Courier{courierMock}
			assignedCourierMock := courierMock // Mock assigned courier

			mockOrderRepo.EXPECT().
				GetFirstInCreatedStatus(ctx).
				Return(orderMock, nil)

			mockCourierRepo.EXPECT().
				GetAllInFreeStatus(ctx).
				Return(couriersMock, nil)

			mockUnitOfWork.EXPECT().
				Begin(ctx).
				Return(ctx)

			mockOrderRepo.EXPECT().
				Update(ctx, orderMock).
				Return(nil)

			mockCourierRepo.EXPECT().
				Update(ctx, assignedCourierMock).
				Return(nil)

			mockUnitOfWork.EXPECT().
				Commit(ctx).
				Return(nil)

			cmd, err := NewAssignOrder()
			require.NoError(t, err)

			assert.NoError(t, handler.Handle(ctx, cmd))
		},
	)
}
