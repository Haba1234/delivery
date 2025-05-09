package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Haba1234/delivery/cmd"
	httpctrl "github.com/Haba1234/delivery/internal/adapters/in/http"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/courier"
	"github.com/Haba1234/delivery/internal/adapters/out/postgres/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/pkg/servers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
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

	configs := cmd.Configs{
		GeoClientURL:                  goDotEnvVariable("GEO_SERVICE_GRPC_URL"),
		KafkaHost:                     goDotEnvVariable("KAFKA_HOST"),
		ConsumerGroup:                 goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaOrdersCreateTopic:        goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
		KafkaOrdersStatusChangedTopic: goDotEnvVariable("KAFKA_DELIVERY_ORDER_CHANGED_TOPIC"),
	}

	connectionString, err := makeConnectionString(dbHost, dbPort, dbUser, dbPassword, dbName, dbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	createDatabaseIfNotExists(dbHost, dbPort, dbUser, dbPassword, dbName)
	gormDB := mustGormOpen(connectionString)
	mustAutoMigrate(gormDB)

	compositionRoot := cmd.NewCompositionRoot(ctx, gormDB, configs)

	cronJob := startCron(compositionRoot)
	defer stopCron(cronJob)

	startKafkaConsumer(ctx, compositionRoot)
	startWebServer(compositionRoot, httpPort)
}

func startKafkaConsumer(ctx context.Context, compositionRoot cmd.CompositionRoot) {
	go func() {
		if err := compositionRoot.Consumers.OrdersCreateConsumer.Consume(ctx); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()
}

func startCron(compositionRoot cmd.CompositionRoot) *cron.Cron {
	c := cron.New()

	_, err := c.AddFunc("@every 1s", compositionRoot.Jobs.AssignOrders.Run)
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}

	_, err = c.AddFunc("@every 2s", compositionRoot.Jobs.MoveCouriers.Run)
	if err != nil {
		log.Fatalf("ошибка при добавлении задачи: %v", err)
	}

	c.Start()

	return c
}

func stopCron(c *cron.Cron) {
	c.Stop()
}

func startWebServer(compositionRoot cmd.CompositionRoot, port string) {
	handlers, err := httpctrl.NewServer(
		compositionRoot.CommandHandlers.CreateOrderHandler,
		compositionRoot.QueryHandlers.GetAllCouriersHandler,
		compositionRoot.QueryHandlers.GetNotCompletedOrdersHandler,
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации HTTP Server: %v", err)
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			},
		),
	)
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET(
		"/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "Healthy")
		},
	)
	registerSwaggerOpenApi(e)
	registerSwaggerUI(e)
	servers.RegisterHandlers(e, handlers)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET(
		"/openapi.json", func(c echo.Context) error {
			swagger, err := servers.GetSwagger()
			if err != nil {
				return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
			}

			data, err := swagger.MarshalJSON()
			if err != nil {
				return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
			}

			return c.Blob(http.StatusOK, "application/json", data)
		},
	)
}

func registerSwaggerUI(e *echo.Echo) {
	e.GET(
		"/docs", func(c echo.Context) error {
			html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
			return c.HTML(http.StatusOK, html)
		},
	)
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func createDatabaseIfNotExists(host, port, user, password, dbName string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", "disable")
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
	err := db.AutoMigrate(&courier.ModelCourier{}, &courier.ModelTransport{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&order.ModelOrder{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}
}
