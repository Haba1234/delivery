package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (*Server) TestCheckout(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
