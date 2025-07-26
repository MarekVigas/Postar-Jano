package api_test

import (
	"context"
	"github.com/MarekVigas/Postar-Jano/internal/services/auth"
	"github.com/MarekVigas/Postar-Jano/internal/services/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/internal/services/promo"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PromoSuite struct {
	CommonSuite
	promoSetup func(commonSuite *CommonSuite) *promo.Registry
}

func NewPromoSuite(promoSetup func(commonSuite *CommonSuite) *promo.Registry) *PromoSuite {
	return &PromoSuite{
		CommonSuite: CommonSuite{},
		promoSetup:  promoSetup,
	}
}

func (s *PromoSuite) SetupSuite() {
	s.CommonSuite.SetupSuite()
	s.promoRegistry = s.promoSetup(&s.CommonSuite)
}

func (s *PromoSuite) TestGeneratePromoCode_Unauthorized() {
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes", echo.Map{
		"email":              "test@example.com",
		"registration_count": 1,
	})
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *PromoSuite) TestGeneratePromoCode_UnprocessableEntity() {
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes", echo.Map{})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{"errors": map[string]interface{}{
			"email":              "invalid",
			"registration_count": "missing",
		}}, body)
	})
}

func (s *PromoSuite) TestGeneratePromoCode_OK_WithoutEmail() {
	const (
		mail              = "test@example.com"
		registrationCount = 1
	)
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes", echo.Map{
		"email":              mail,
		"registration_count": registrationCount,
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.NotEmpty(body["promo_code"])
	})
}

func (s *PromoSuite) TestGeneratePromoCode_OK() {
	const (
		mail              = "test@example.com"
		registrationCount = 1
	)

	s.mailer.On("PromoMail", mock.Anything, mock.MatchedBy(
		func(req *templates.PromoReq) bool {
			return req.Mail == mail &&
				req.AvailableRegistrations == registrationCount
		})).Return(nil)

	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes", echo.Map{
		"email":              mail,
		"registration_count": registrationCount,
		"send_email":         true,
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Require().NotEmpty(body["promo_code"])
		promoCode, ok := body["promo_code"].(string)
		s.Require().True(ok)
		dbPromoCode, err := s.promoRegistry.ValidateTokenWithQueryerContext(context.Background(), s.dbx, promoCode)
		s.Require().NoError(err)
		s.Equal(dbPromoCode.Email, mail)
	})

	s.mailer.AssertExpectations(s.T())
}

func (s *PromoSuite) TestValidatePromoCode_UnprocessableEntity() {
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes/validate", echo.Map{})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{"errors": map[string]interface{}{
			"promo_code": "missing",
		}}, body)
	})
}

func (s *PromoSuite) TestValidatePromoCode_NonExisting() {
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes/validate", echo.Map{
		"promo_code": "XYZ",
	})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"status": "invalid", "available_registrations": float64(0)}, body)
	})
}

func (s *PromoSuite) TestValidatePromoCode_AlreadyUsed() {
	ctx := context.Background()
	token, err := s.promoRegistry.GenerateToken(ctx, "test@test.com", 1, false)
	s.Require().NoError(err)

	promoCode, err := s.promoRegistry.ValidateTokenWithQueryerContext(ctx, s.dbx, token)
	s.Require().NoError(err)

	s.Require().NoError(s.promoRegistry.MarkTokenUsage(ctx, s.dbx, promoCode.Key))

	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes/validate", echo.Map{"promo_code": token})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"status": "invalid", "available_registrations": float64(0)}, body)
	})
}

func (s *PromoSuite) TestValidatePromoCode_Valid() {
	const regCount = 10
	token, err := s.promoRegistry.GenerateToken(context.Background(), "test@test.com", regCount, false)
	s.Require().NoError(err)

	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes/validate", echo.Map{"promo_code": token})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"status": "ok", "available_registrations": float64(regCount)}, body)
	})
}

func TestSimplePromoSuite(t *testing.T) {
	suite.Run(t, NewPromoSuite(func(s *CommonSuite) *promo.Registry {
		return promo.NewRegistry(s.postgresDB, promo.NewSimpleGenerator(s.logger), s.mailer)
	}))
}

func TestJWTPromoSuite(t *testing.T) {
	suite.Run(t, NewPromoSuite(func(s *CommonSuite) *promo.Registry {
		return promo.NewRegistry(s.postgresDB, promo.NewJWTGenerator(s.logger, []byte(promoSecret), nil, nil), s.mailer)
	}))
}
