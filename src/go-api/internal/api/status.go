package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *API) Status(c echo.Context) error {
	resp, ok := api.checker.Ping(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusInternalServerError, resp)
	}
	return c.JSON(http.StatusOK, resp)
}
