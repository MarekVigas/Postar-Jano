package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/internal/promo"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PromoSuite struct {
	CommonSuite
	promoSetup func(commonSuite *CommonSuite) repository.PromoManager
}

func NewPromoSuite(promoSetup func(commonSuite *CommonSuite) repository.PromoManager) *PromoSuite {
	return &PromoSuite{
		CommonSuite: CommonSuite{},
		promoSetup:  promoSetup,
	}
}

func (s *PromoSuite) SetupSuite() {
	s.CommonSuite.SetupSuite()
	s.promoSetup(&s.CommonSuite)
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
		dbPromoCode, err := s.promoManager.ValidateToken(context.Background(), s.dbx, promoCode)
		s.Require().NoError(err)
		s.Equal(dbPromoCode.Email, mail)
	})

	s.mailer.AssertExpectations(s.T())
}

func (s *PromoSuite) TestValidatePromoCode_UnprocessableEntity() {

}

func (s *PromoSuite) TestValidatePromoCode_NonExisting() {

}

func (s *PromoSuite) TestValidatePromoCode_AlreadyUsed() {

}

func (s *PromoSuite) TestValidatePromoCode_Valid() {

}

func TestSimplePromoSuite(t *testing.T) {
	suite.Run(t, NewPromoSuite(func(s *CommonSuite) repository.PromoManager {
		return promo.NewSimpleGenerator(s.logger)
	}))
}

func TestJWTPromoSuite(t *testing.T) {
	suite.Run(t, NewPromoSuite(func(s *CommonSuite) repository.PromoManager {
		return promo.NewJWTGenerator(s.logger, []byte(promoSecret), nil, nil)
	}))
}
