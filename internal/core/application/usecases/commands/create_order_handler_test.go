package commands

import (
	"context"
	"testing"

	orderrepo "github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleWithMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mocks.NewMockIOrderRepository(ctrl)
	mockGeoClient := mocks.NewMockIGeoClient(ctrl)

	commandHandler, err := NewCreateOrderHandler(mockOrderRepo, mockGeoClient)
	require.NoError(t, err)

	t.Run(
		"empty command", func(t *testing.T) {
			cmd := CreateOrder{}
			err := commandHandler.Handle(t.Context(), cmd)
			require.Error(t, err)
			assert.EqualError(t, err, "value is required: add address command")
		},
	)

	t.Run(
		"order already exists", func(t *testing.T) {
			orderID := uuid.New()
			existingOrder := &order.Order{}

			mockOrderRepo.EXPECT().
				Get(gomock.Any(), orderID).
				Return(existingOrder, nil)

			cmd, err := NewCreateOrder(orderID, "Новая улица")
			require.NoError(t, err)

			err = commandHandler.Handle(t.Context(), cmd)
			assert.ErrorIs(t, err, ErrOrderAlreadyExists)
		},
	)

	t.Run(
		"new order was created", func(t *testing.T) {
			orderID := uuid.New()
			location, err := kernel.CreateRandomLocation()
			require.NoError(t, err)

			mockOrderRepo.EXPECT().
				Get(gomock.Any(), orderID).
				Return(nil, orderrepo.ErrOrderNotFound)

			mockGeoClient.EXPECT().
				GetGeolocation(gomock.Any(), "Новая улица").
				Return(location, nil)

			mockOrderRepo.EXPECT().
				Add(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(_ context.Context, o *order.Order) error {
						assert.Equal(t, o.ID(), orderID)
						assert.False(t, o.Location().IsEmpty())
						return nil
					},
				)

			cmd, err := NewCreateOrder(orderID, "Новая улица")
			require.NoError(t, err)

			assert.NoError(t, commandHandler.Handle(t.Context(), cmd))
		},
	)
}
