package api

import (
	"context"
	"database/sql"
	"github.com/MarekVigas/Postar-Jano/internal/resources"
	"github.com/MarekVigas/Postar-Jano/internal/services/events"
	"github.com/MarekVigas/Postar-Jano/internal/services/promo"
	"github.com/MarekVigas/Postar-Jano/internal/services/registration"
	"github.com/MarekVigas/Postar-Jano/internal/services/status"
	"github.com/MarekVigas/Postar-Jano/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(ctx context.Context, username string, pass string) (string, error)
	Middleware() echo.MiddlewareFunc
}

type API struct {
	*echo.Echo
	authenticator       Authenticator
	eventManager        *events.Manager
	promoRegistry       *promo.Registry
	registrationManager *registration.Manager
	checker             *status.Checker
	cookieDomain        string
	cookieSecure        bool
}

func New(
	log *zap.Logger,
	authenticator Authenticator,
	eventManager *events.Manager,
	promoRegistry *promo.Registry,
	registrationManager *registration.Manager,
	checker *status.Checker,
	adminOrigin string,
	cookieDomain string,
	cookieSecure bool,
) *API {
	e := echo.New()
	a := &API{
		Echo:                e,
		authenticator:       authenticator,
		eventManager:        eventManager,
		promoRegistry:       promoRegistry,
		registrationManager: registrationManager,
		checker:             checker,
		cookieDomain:        cookieDomain,
		cookieSecure:        cookieSecure,
	}

	requireAuth := authenticator.Middleware()
	registry := prometheus.NewRegistry()
	e.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{adminOrigin, "https://leto.salezko.sk"},
			AllowCredentials: true,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		}),
		middleware.RequestID(),
		logger.ContextLogger(log),
		echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
			Subsystem:  "postar_jano",
			Registerer: registry,
		}),
	)
	e.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: registry,
	}))
	api := e.Group("/api",
		middleware.Recover(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().URL.Path == "/api/status" {
					return true
				}
				if strings.HasPrefix(c.Request().URL.Path, "/api/stats") {
					return true
				}
				return false
			},
		}),
	)
	api.GET("/status", a.Status)

	api.POST("/registrations/:id", a.Register)

	api.GET("/stats", a.ListStats)
	api.GET("/stats/:id", a.StatByID)

	api.GET("/events", a.ListEvents)
	api.GET("/events/:id", a.EventByID)

	api.GET("/registrations/:token", a.FindRegistration)
	api.POST("/promo_codes/validate", a.ValidatePromoCode)

	// Admin
	api.POST("/sign/in", a.SignIn)
	api.POST("/sign/out", a.SignOut)
	api.GET("/me", a.Me, requireAuth)
	api.POST("/promo_codes", a.GeneratePromoCode, requireAuth)
	api.GET("/registrations", a.ListRegistrations, requireAuth)
	api.GET("/registrations/:id", a.FindRegistrationByID, requireAuth)
	api.DELETE("/registrations/:id", a.DeleteRegistrationByID, requireAuth)
	api.PUT("/registrations/:id", a.UpdateRegistration, requireAuth)
	api.POST("/registrations/:id/resend_confirmation", a.ResendConfirmation, requireAuth)
	api.POST("/send_payment_notifications", a.SendPaymentNotification, requireAuth)
	api.POST("/send_payment_reminder", a.SendPaymentReminder, requireAuth)

	return a
}

func (api *API) SignIn(c echo.Context) error {
	var req resources.SignIn
	if err := c.Bind(&req); err != nil {
		return err
	}

	if errs := req.Validate(); errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, errs)
	}

	ctx := c.Request().Context()
	token, err := api.authenticator.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		switch err := errors.Cause(err); {
		case errors.Is(err, sql.ErrNoRows):
			return echo.ErrForbidden
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return echo.ErrForbidden
		default:
			logger.FromCtx(ctx).Error("Error during authentication.", zap.Error(err), zap.String("username", req.Username))
			return echo.ErrForbidden
		}
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   api.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   3 * 60 * 60,
	}
	if api.cookieDomain != "" {
		cookie.Domain = api.cookieDomain
	}
	c.SetCookie(cookie)
	return c.NoContent(http.StatusNoContent)
}

func (api *API) SignOut(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   api.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	}
	if api.cookieDomain != "" {
		cookie.Domain = api.cookieDomain
	}
	c.SetCookie(cookie)
	return c.NoContent(http.StatusNoContent)
}

func (api *API) Me(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return echo.ErrUnauthorized
	}
	email, _ := claims["sub"].(string)
	return c.JSON(http.StatusOK, echo.Map{"email": email})
}

func (api *API) handleError(err error) error {
	switch err := errors.Cause(err); {
	case errors.Is(err, sql.ErrNoRows):
		return echo.ErrNotFound
	default:
		return err
	}
}

func (api *API) getIntParam(c echo.Context, name string) (int, error) {
	id, err := strconv.ParseInt(c.Param(name), 10, 32)
	if err != nil {
		return 0, echo.ErrBadRequest
	}
	return int(id), nil
}

func validateStruct(s interface{}) interface{} {
	v := validator.New()
	err := v.Struct(s)
	if err == nil {
		return nil
	}
	var validationErrs validator.ValidationErrors
	if ok := errors.As(err, &validationErrs); !ok {
		return err
	}

	t := reflect.TypeOf(s).Elem()
	mustFieldName := func(fieldName string) string {
		tokens := strings.Split(fieldName, ".")
		var res []string
		current := t
		for _, token := range tokens[1:] {
			field, ok := current.FieldByName(token)
			if !ok {
				panic("field not found:" + token)
			}
			jsonName := field.Tag.Get("json")
			if jsonName == "" {
				res = append(res, token)
			}
			res = append(res, jsonName)
			current = field.Type
		}

		return strings.Join(res, ".")
	}

	errs := echo.Map{}
	for _, fieldErr := range validationErrs {
		var errName string
		switch fieldErr.Tag() {
		case "required":
			errName = "missing"
		default:
			errName = "invalid"
		}
		errs[mustFieldName(fieldErr.Namespace())] = errName
	}
	return errs
}
