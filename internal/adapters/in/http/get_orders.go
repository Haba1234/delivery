package http

import (
	"errors"
	"net/http"

	"github.com/Haba1234/delivery/internal/adapters/in/http/problems"
	"github.com/Haba1234/delivery/internal/core/application/usecases/queries"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/pkg/servers"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetOrders(c echo.Context) error {
	query, err := queries.NewGetNotCompletedOrders()
	if err != nil {
		return problems.NewBadRequestError(err.Error())
	}

	response, err := s.getNotCompletedOrdersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFoundError(err.Error()))
		}
	}

	var orders []servers.Order
	for _, o := range response.Orders {
		location := servers.Location{
			X: o.LocationX,
			Y: o.LocationY,
		}

		var order = servers.Order{
			Id:       o.ID,
			Location: location,
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}
