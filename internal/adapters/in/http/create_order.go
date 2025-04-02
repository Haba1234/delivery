package http

import (
	"math/rand"
	"net/http"

	"github.com/Haba1234/delivery/internal/adapters/in/http/problems"
	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Server) CreateOrder(c echo.Context) error {
	createOrderCommand, err := commands.NewCreateOrder(uuid.New(), randomStreetName())
	if err != nil {
		return problems.NewBadRequestError(err.Error())
	}

	err = s.createOrderCommandHandler.Handle(c.Request().Context(), createOrderCommand)
	if err != nil {
		return problems.NewConflictError(err.Error(), "/")
	}

	return c.JSON(http.StatusOK, nil)
}

func randomStreetName() string {
	streetNames := []string{
		"Тестировочная", "Айтишная", "Эйчарная", "Аналитическая", "Нагрузочная", "Серверная", "Мобильная", "Бажная",
		"Несуществующая",
	}
	return streetNames[rand.Intn(len(streetNames))]
}
