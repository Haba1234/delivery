package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartPostgresContainer(ctx context.Context) (testcontainers.Container, string, error) {
	dbName := "testdb"
	user := "testuser"
	password := "testpass"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(
		ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		return nil, "", err
	}

	// Получаем параметры подключения
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, "", err
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, "", err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port.Port(), dbName)

	return postgresContainer, dsn, nil
}
