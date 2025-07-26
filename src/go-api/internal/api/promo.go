package api

import (
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (api *API) GeneratePromoCode(c echo.Context) error {
	// Decode request
	var req resources.PromoCodeReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	if errs := validateStruct(&req); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"errors": errs})
	}

	ctx := c.Request().Context()
	token, err := api.promoRegistry.GenerateToken(ctx, req.Email, req.RegistrationCount, req.SendEmail)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, echo.Map{"promo_code": token})
}

func (api *API) ValidatePromoCode(c echo.Context) error {
	var req resources.ValidatePromoCodeReq
	if err := c.Bind(&req); err != nil {
		return err
	}
	if errs := validateStruct(&req); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"errors": errs})
	}

	resp, err := api.promoRegistry.ValidateToken(c.Request().Context(), req.PromoCode)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, resp)
}
