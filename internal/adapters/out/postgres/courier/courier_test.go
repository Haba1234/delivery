package courier

import (
	"context"
	"fmt"
	"testing"

	"github.com/Haba1234/delivery/internal/core/domain/model/courier"
	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&ModelCourier{})
	require.NoError(t, err)

	err = db.AutoMigrate(&ModelTransport{})
	require.NoError(t, err)

	return ctx, db, nil
}

func generateTestCouriers(t *testing.T) []*courier.Courier {
	t.Helper()

	var couriers []*courier.Courier

	for i := 0; i < 5; i++ {
		location, err := kernel.CreateRandomLocation() // Генерация случайной локации
		require.NoError(t, err)

		courierAggregate, err := courier.New(
			fmt.Sprintf("Курьер_%v", i),
			fmt.Sprintf("Транспорт_%v", i),
			2,
			location,
		)
		require.NoError(t, err)

		couriers = append(couriers, courierAggregate)
	}

	return couriers
}

func TestCourierRepository(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	require.NoError(t, err)

	// Создаем репозиторий
	courierRepository, err := NewRepository(db)
	require.NoError(t, err)

	// Подготовка данных
	couriers := generateTestCouriers(t)

	// Добавляем курьеров в БД
	for _, c := range couriers {
		err = courierRepository.Add(ctx, c)
		require.NoError(t, err)
	}

	// Проверяем первого курьера
	firstCourier, err := courierRepository.Get(ctx, couriers[0].ID())
	require.NoError(t, err)
	assert.Equal(t, couriers[0], firstCourier)

	// Получаем всех свободных курьеров
	freeCouriers, err := courierRepository.GetAllInFreeStatus(ctx)
	require.NoError(t, err)
	assert.ElementsMatch(t, couriers, freeCouriers)
}
