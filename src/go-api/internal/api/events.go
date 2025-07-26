package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *API) ListEvents(c echo.Context) error {
	events, err := api.eventManager.GetAll(c.Request().Context())
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, events)
}

func (api *API) EventByID(c echo.Context) error {

	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	event, err := api.eventManager.GetByID(c.Request().Context(), id)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, event)
}

func (api *API) ListStats(c echo.Context) error {
	stats, err := api.eventManager.GetAllStats(c.Request().Context())
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, stats)
}

func (api *API) StatByID(c echo.Context) error {
	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	stats, err := api.eventManager.GetStatById(c.Request().Context(), id)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, stats)
}
