package order

import (
	"context"
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/pkg/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTest инициализирует окружение для тестирования.
func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	t.Helper()

	ctx := t.Context()

	postgresContainer, dsn, err := testutil.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Очистка выполняется после завершения теста
	t.Cleanup(
		func() {
			assert.NoError(t, postgresContainer.Terminate(context.Background()))
		},
	)

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Авто миграция (создаём таблицу для заказов)
	err = db.AutoMigrate(&ModelOrder{})
	require.NoError(t, err)

	return ctx, db, nil
}

// generateTestOrders создает тестовые данные для заказов.
func generateTestOrders(t *testing.T) []*order.Order {
	t.Helper()

	var orders []*order.Order

	for i := 0; i < 5; i++ {
		location, err := kernel.CreateRandomLocation() // Генерация случайной локации
		require.NoError(t, err)

		orderAggregate, err := order.New(uuid.New(), location)
		require.NoError(t, err)

		orders = append(orders, orderAggregate)
	}

	return orders
}

func TestOrderRepository(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	require.NoError(t, err)

	// Создаем репозиторий
	orderRepository, err := NewRepository(db)
	require.NoError(t, err)

	// Подготовка данных
	testOrders := generateTestOrders(t)

	// Добавляем заказы в БД
	for _, o := range testOrders {
		err = orderRepository.Add(ctx, o)
		require.NoError(t, err)
	}

	// Проверяем первый заказ
	firstOrder, err := orderRepository.Get(ctx, testOrders[0].ID())
	require.NoError(t, err)
	assert.Equal(t, testOrders[0], firstOrder)

	// Получаем все созданные заказы
	firstOrder, err = orderRepository.GetFirstInCreatedStatus(ctx)
	require.NoError(t, err)
	assert.Contains(t, testOrders, firstOrder)

	// На третий заказ назначаем курьера
	courierLocation, err := kernel.CreateRandomLocation()
	require.NoError(t, err)

	couirerJon, err := courier.New("Jon", "Bike", 2, courierLocation)
	require.NoError(t, err)

	require.NoError(t, testOrders[3].Assign(couirerJon))
	require.NoError(t, orderRepository.Update(ctx, testOrders[3]))

	// Получаем все назначенные заказы
	assignedOrders, err := orderRepository.GetAllInAssignedStatus(ctx)
	require.NoError(t, err)
	assert.Subset(t, testOrders, assignedOrders)
}
