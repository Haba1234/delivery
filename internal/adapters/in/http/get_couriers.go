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

func (s *Server) GetCouriers(c echo.Context) error {
	query, err := queries.NewGetAllCouriers()
	if err != nil {
		return problems.NewBadRequestError(err.Error())
	}

	response, err := s.getAllCouriersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFoundError(err.Error()))
		}
	}

	var couriers []servers.Courier
	for _, courier := range response.Couriers {
		location := servers.Location{
			X: courier.LocationX,
			Y: courier.LocationY,
		}

		var courier = servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}
		couriers = append(couriers, courier)
	}
	return c.JSON(http.StatusOK, couriers)
}
