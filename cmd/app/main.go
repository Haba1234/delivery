package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Haba1234/delivery/cmd"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()
	httpPort := goDotEnvVariable("HTTP_PORT")
	dbHost := goDotEnvVariable("DB_HOST")
	dbPort := goDotEnvVariable("DB_PORT")
	dbUser := goDotEnvVariable("DB_USER")
	dbPassword := goDotEnvVariable("DB_PASSWORD")
	dbName := goDotEnvVariable("DB_DBNAME")
	dbSslMode := goDotEnvVariable("DB_SSLMODE")
	connectionString, err := makeConnectionString(dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	createDatabaseIfNotExists(dbHost, dbPort, dbName)
	gormDB := mustGormOpen(connectionString)
	mustAutoMigrate(gormDB)

	compositionRoot := cmd.NewCompositionRoot(
		ctx,
		gormDB,
	)
	startWebServer(compositionRoot, httpPort)
}

func startWebServer(compositionRoot cmd.CompositionRoot, port string) {
	e := echo.New()
	e.GET(
		"/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "Healthy")
		},
	)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func createDatabaseIfNotExists(host, port, dbName string) {
	dsn, err := makeConnectionString(host, port, "postgres", "postgres", "postgres", "disable")
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Создаём базу данных, если её нет
	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		log.Printf("Ошибка создания БД (возможно, уже существует): %v", err)
	}
}

func makeConnectionString(host, port, user, password, dbName, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError(host)
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError(port)
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError(user)
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError(password)
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError(dbName)
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError(sslMode)
	}
	return fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	), nil
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  connectionString,
				PreferSimpleProtocol: true,
			},
		), &gorm.Config{},
	)
	if err != nil {
		log.Fatalf("connection to postgres through gorm\n: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&courier.ModelCourier{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&courier.ModelCourier{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&order.ModelOrder{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
}
