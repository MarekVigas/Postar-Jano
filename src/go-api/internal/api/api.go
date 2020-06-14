package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/internal/resources"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const tokenLifetime = 3 * time.Hour

type EmailSender interface {
	ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error
}

type Authenticator interface {
	Authenticate(ctx context.Context, username string, pass string) (*model.Owner, error)
}

type API struct {
	*echo.Echo
	repo          *repository.PostgresRepo
	logger        *zap.Logger
	authenticator Authenticator
	sender        EmailSender
	jwtSecret     []byte
}

func New(
	logger *zap.Logger,
	repo *repository.PostgresRepo,
	authenticator Authenticator,
	sender EmailSender,
	jwtSecret []byte,
) *API {
	e := echo.New()
	a := &API{
		Echo:          e,
		repo:          repo,
		logger:        logger,
		authenticator: authenticator,
		sender:        sender,
		jwtSecret:     jwtSecret,
	}

	jwt := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.Claims{},
		SigningKey: jwtSecret,
		ErrorHandler: func(err error) error {
			return echo.ErrUnauthorized
		},
	})

	api := e.Group("/api",
		middleware.Recover(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Skipper: func(c echo.Context) bool {
				if c.Request().URL.Path == "/api/stats" {
					return true
				}
				return false
			},
		}))
	api.GET("/status", a.Status)

	api.POST("/registrations/:id", a.Register)

	api.GET("/stats", a.ListStats)
	api.GET("/stats/:id", a.StatByID)

	api.GET("/events", a.ListEvents)
	api.GET("/events/:id", a.EventByID)

	// TODO: will be delivered after Sunday :)
	api.GET("/registrations/:token", a.FindRegistration)

	api.POST("/sign/in", a.SignIn)
	api.GET("/registrations", a.ListRegistrations, jwt)
	api.PUT("/registrations/:id", a.UpdateRegistration, jwt)

	return a
}

func (api *API) Status(c echo.Context) error {
	if err := api.repo.Ping(c.Request().Context()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"status": "err"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (api *API) ListStats(c echo.Context) error {
	return nil
}

func (api *API) StatByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	stats, err := api.repo.GetStat(ctx, id)
	if err != nil {
		api.Logger.Error("failed to list stats", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (api *API) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var req resources.RegisterReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	api.logger.Debug("Request received", zap.Reflect("raw", req))
	// TODO: Validate input
	if len(req.DayIDs) == 0 {
		return c.JSON(http.StatusUnprocessableEntity, "No days provided.")
	}

	eventID, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	reg, ok, err := api.repo.Register(ctx, &req, eventID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return echo.ErrNotFound
		}
		api.logger.Error("Failed to create a registration", zap.Error(err))
		return err
	}

	if !ok {
		return c.JSON(http.StatusOK, echo.Map{
			"success":       reg.Success,
			"registeredIDs": reg.RegisteredIDs,
		})
	}

	pills := "-"
	if reg.Reg.Pills != nil {
		pills = *reg.Reg.Pills
	}

	restrictions := "-"
	if reg.Reg.Problems != nil {
		restrictions = *reg.Reg.Problems
	}
	var info string
	if reg.Event.MailInfo != nil {
		info = *reg.Event.MailInfo
	}

	// Send confirmation mail.
	if err := api.sender.ConfirmationMail(ctx, &templates.ConfirmationReq{
		Mail:          reg.Reg.Email,
		ParentName:    reg.Reg.ParentName,
		ParentSurname: reg.Reg.ParentSurname,
		EventName:     reg.Event.Title,
		Name:          reg.Reg.Name,
		Surname:       reg.Reg.Surname,
		Pills:         pills,
		Restrictions:  restrictions,
		Text:          reg.Event.OwnerPhone + " " + reg.Event.OwnerEmail,
		PhotoURL:      reg.Event.OwnerPhoto,
		Sum:           reg.Reg.Amount,
		Owner:         reg.Event.OwnerName + " " + reg.Event.OwnerSurname,
		Days:          reg.RegisteredDesc,
		Info:          info,
	}); err != nil {
		api.logger.Error("Failed to send a confirmation mail.", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success":       reg.Success,
		"registeredIDs": reg.RegisteredIDs,
		"token":         reg.Token,
	})
}

func (api *API) ListRegistrations(c echo.Context) error {
	ctx := c.Request().Context()
	regs, err := api.repo.ListRegistrations(ctx)
	if err != nil {
		api.logger.Error("Failed to list registrations", zap.Error(err))
		return err
	}
	if len(regs) == 0 {
		regs = []model.ExtendedRegistration{}
	}
	return c.JSON(http.StatusOK, regs)
}

func (api *API) FindRegistration(c echo.Context) error {
	return nil
}

func (api *API) ListEvents(c echo.Context) error {
	ctx := c.Request().Context()

	events, err := api.repo.ListEvents(ctx)
	if err != nil {
		api.Logger.Error("Failed to list events.", zap.Error(err))
		return err
	}
	return c.JSON(http.StatusOK, events)
}

func (api *API) EventByID(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := api.getIntParam(c, "id")
	if err != nil {
		return err
	}

	event, err := api.repo.FindEvent(ctx, id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return echo.ErrNotFound
		}
		api.Logger.Error("Failed to find event.", zap.Error(err), zap.Int("event_id", id))
		return err
	}
	return c.JSON(http.StatusOK, event)
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
	owner, err := api.authenticator.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		switch errors.Cause(err) {
		case sql.ErrNoRows:
			return echo.ErrForbidden
		case bcrypt.ErrMismatchedHashAndPassword:
			return echo.ErrForbidden
		default:
			api.logger.Error("Error during authentication.", zap.Error(err), zap.String("username", req.Username))
			return echo.ErrForbidden
		}
	}
	token, err := api.generateToken(owner)
	if err != nil {
		api.logger.Error("Failed to generate token.", zap.Error(err))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
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
	if err := api.repo.UpdateRegistrations(ctx, &model.Registration{
		ID:        id,
		Name:      req.Child.Name,
		Surname:   req.Child.Surname,
		Email:     req.Parent.Email,
		Payed:     req.Payed,
		Discount:  req.Discount,
		AdminNote: req.AdminNote,
	}); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return echo.ErrNotFound
		}
		api.logger.Error("Failed to update registration.", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusAccepted, nil)
}

func (api *API) generateToken(owner *model.Owner) (string, error) {
	now := time.Now()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Add(tokenLifetime).Unix(),
			Id:        owner.Email,
			IssuedAt:  now.Unix(),
			Issuer:    "",
			NotBefore: 0,
			Subject:   "",
		},
	})

	return tok.SignedString(api.jwtSecret)
}

func (api *API) getIntParam(c echo.Context, name string) (int, error) {
	id, err := strconv.ParseInt(c.Param(name), 10, 32)
	if err != nil {
		return 0, echo.ErrBadRequest
	}
	return int(id), nil
}
