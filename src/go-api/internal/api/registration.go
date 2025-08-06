package api

import (
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/MarekVigas/Postar-Jano/internal/services/promo"
	"github.com/MarekVigas/Postar-Jano/internal/services/registration"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

func (api *API) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var req resources.RegisterReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	logger.FromCtx(ctx).Debug("Request received", zap.Reflect("raw", req))

	if errs := validateStruct(&req); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"errors": errs})
	}

	eventID, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	reg, err := api.registrationManager.CreateNew(ctx, &req, eventID)
	if err != nil {
		switch errors.Cause(err) {
		case registration.ErrNotActive:
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"errors": map[string]interface{}{
					"event_id": "not active",
				}})
		case promo.ErrAlreadyUsed:
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Token already used.",
			})
		default:
			return api.handleError(err)
		}
	}

	return c.JSON(http.StatusOK, reg)
}

func (api *API) ListRegistrations(c echo.Context) error {
	ctx := c.Request().Context()
	regs, err := api.registrationManager.GetAll(ctx)
	if err != nil {
		return api.handleError(err)
	}
	if len(regs) == 0 {
		regs = []model.ExtendedRegistration{}
	}
	return c.JSON(http.StatusOK, regs)
}

func (api *API) FindRegistration(c echo.Context) error {
	ctx := c.Request().Context()
	reg, err := api.registrationManager.GetByToken(ctx, c.Param("token"))
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, reg)
}

func (api *API) FindRegistrationByID(c echo.Context) error {
	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}
	reg, err := api.registrationManager.GetByID(c.Request().Context(), id)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, reg)
}

func (api *API) DeleteRegistrationByID(c echo.Context) error {
	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}
	reg, err := api.registrationManager.DeleteByID(c.Request().Context(), id)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, reg)
}

func (api *API) UpdateRegistration(c echo.Context) error {
	var req resources.UpdateReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	if errs := req.Validate(); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	ctx := c.Request().Context()
	if err := api.registrationManager.Update(ctx, &model.Registration{
		ID:        id,
		Amount:    req.Amount,
		Payed:     req.Payed,
		AdminNote: req.AdminNote,
	}); err != nil {
		return api.handleError(err)
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (api *API) ResendConfirmation(c echo.Context) error {
	regID, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	var req resources.ResendConfirmationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if errs := validateStruct(&req); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"errors": errs})
	}

	if err := api.registrationManager.ResendConfirmation(c.Request().Context(), regID, req.Email); err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusAccepted, nil)

}

func (api *API) SendPaymentNotification(c echo.Context) error {
	resp, err := api.registrationManager.SendPaymentNotifications(c.Request().Context())
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (api *API) SendPaymentReminder(c echo.Context) error {
	var req resources.PaymentReminderRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := api.registrationManager.SendPaymentReminder(c.Request().Context(), req.EventId)
	if err != nil {
		return api.handleError(err)
	}
	return c.JSON(http.StatusOK, resp)
}
