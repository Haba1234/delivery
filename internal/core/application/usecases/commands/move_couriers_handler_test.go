package commands

import (
	"context"
	"testing"

	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMoveCouriersHandler_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockIOrderRepository(ctrl)
	mockCourierRepo := mocks.NewMockICourierRepository(ctrl)
	mockUnitOfWork := mocks.NewMockIUnitOfWork(ctrl)

	handler, err := NewMoveCouriersHandler(mockUnitOfWork, mockOrderRepo, mockCourierRepo)
	require.NoError(t, err)

	t.Run(
		"empty command", func(t *testing.T) {
			cmd := MoveCouriers{}
			err := handler.Handle(t.Context(), cmd)
			require.Error(t, err)
			assert.EqualError(t, err, "value is required add address command")
		},
	)

	t.Run(
		"assigned orders not found", func(t *testing.T) {
			mockOrderRepo.EXPECT().
				GetAllInAssignedStatus(gomock.Any()).
				Return(nil, orderrepo.ErrAssignedOrdersNotFound)

			cmd, err := NewMoveCouriers()
			require.NoError(t, err)

			err = handler.Handle(t.Context(), cmd)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"successful move and completion", func(t *testing.T) {
			ctx := t.Context()

			orderLocation, err := kernel.CreateLocation(2, 1)
			require.NoError(t, err)
			orderMock, err := order.New(uuid.New(), orderLocation)
			require.NoError(t, err)

			courierLocation, err := kernel.CreateLocation(1, 1)
			require.NoError(t, err)

			courierMock, err := courier.New("Jon", "bike", 1, courierLocation)

			require.NoError(t, orderMock.Assign(courierMock))

			ordersMock := []*order.Order{orderMock}

			expectedLocation := orderMock.Location()

			mockOrderRepo.EXPECT().
				GetAllInAssignedStatus(gomock.Any()).
				Return(ordersMock, nil)

			mockCourierRepo.EXPECT().
				Get(gomock.Any(), *orderMock.CourierID()).
				Return(courierMock, nil)

			mockUnitOfWork.EXPECT().
				Begin(ctx).
				Return(ctx)

			mockCourierRepo.EXPECT().
				Update(ctx, gomock.Any()).
				DoAndReturn(
					func(_ context.Context, c *courier.Courier) error {
						assert.True(t, c.Location().Equals(expectedLocation))
						return nil
					},
				)

			mockOrderRepo.EXPECT().
				Update(ctx, gomock.Any()).
				DoAndReturn(
					func(_ context.Context, o *order.Order) error {
						assert.True(t, o.IsCompleted())
						return nil
					},
				)

			mockUnitOfWork.EXPECT().
				Commit(ctx).
				Return(nil)

			cmd, err := NewMoveCouriers()
			require.NoError(t, err)

			assert.NoError(t, handler.Handle(ctx, cmd))
		},
	)
}
