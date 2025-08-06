package api_test

import (
	"context"
	"fmt"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"
	"github.com/MarekVigas/Postar-Jano/internal/services/auth"
	"github.com/MarekVigas/Postar-Jano/internal/services/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/internal/services/promo"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RegistrationSuite struct {
	CommonSuite
}

func (s *RegistrationSuite) TestRegister_UnprocessableEntity() {
	event := s.InsertEvent()

	u := fmt.Sprintf("/api/registrations/%d", event.ID)
	req, rec := s.NewRequest(http.MethodPost, u, nil)

	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{
				"child.name":               "missing",
				"child.surname":            "missing",
				"child.city":               "missing",
				"child.dateOfBirth":        "missing",
				"child.finishedSchoolYear": "missing",
				"child.gender":             "invalid",
				"parent.email":             "invalid",
				"parent.name":              "missing",
				"parent.surname":           "missing",
				"parent.phone":             "missing",
				"days":                     "missing",
			},
		}, body)
	})
}

func (s *RegistrationSuite) TestRegister_NotFoundDay() {
	//TODO
}

func (s *RegistrationSuite) TestRegister_NotFoundEvent() {
	req, rec := s.NewRequest(http.MethodPost, "/api/registrations/42", echo.Map{
		"child": echo.Map{
			"name":               "meno",
			"surname":            "priezvisko",
			"gender":             "female",
			"city":               "city",
			"finishedSchoolYear": "school",
			"dateOfBirth":        time.Now().Format(time.RFC3339),
		},
		"parent": echo.Map{
			"name":    "pname",
			"surname": "psurname",
			"email":   "email@email.com",
			"phone":   "phone",
		},
		"days": []interface{}{9, 500},
	})

	s.AssertServerResponseObject(req, rec, http.StatusNotFound, func(body echo.Map) {
		s.Equal(body, echo.Map{"message": "Not Found"})
	})
}

func (s *RegistrationSuite) TestRegister_NotActive() {
	const (
		name     = "dano"
		surname  = "zharmanca"
		pname    = "janko"
		psurname = "hrasko"
		gender   = "male"
		city     = "BB"
		phone    = "+421"
		email    = "dano@mail.sk"
		school   = "3.ZS"
	)
	event := s.InsertEvent()

	_, err := s.dbx.Exec("UPDATE events SET active = false WHERE id = $1", event.ID)
	s.Require().NoError(err)

	day := event.Days[0]
	birth := time.Now().Format(time.RFC3339)
	u := fmt.Sprintf("/api/registrations/%d", event.ID)
	req, rec := s.NewRequest(http.MethodPost, u, echo.Map{
		"child": echo.Map{
			"name":               name,
			"surname":            surname,
			"gender":             gender,
			"city":               city,
			"finishedSchoolYear": school,
			"dateOfBirth":        birth,
		},
		"parent": echo.Map{
			"name":    pname,
			"surname": psurname,
			"email":   email,
			"phone":   phone,
		},
		"days": []interface{}{day.ID},
	})

	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{
				"event_id": "not active",
			},
		}, body)
	})
}

func (s *RegistrationSuite) testRegister_NotActive_Promo() {
	const (
		name     = "dano"
		surname  = "zharmanca"
		pname    = "janko"
		psurname = "hrasko"
		gender   = "male"
		city     = "BB"
		phone    = "+421"
		email    = "dano@mail.sk"
		school   = "3.ZS"
	)
	event := s.InsertEvent()

	// Generate promo code
	req, rec := s.NewRequest(http.MethodPost, "/api/promo_codes", echo.Map{
		"email":              "test@example.com",
		"registration_count": 1,
	})
	s.AuthorizeRequest(req, &auth.Claims{})

	var promoCode string
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		val, ok := body["promo_code"]
		s.Require().True(ok)
		promoCode, ok = val.(string)
		s.Require().True(ok)
	})

	day := event.Days[0]
	birth := time.Now().Format(time.RFC3339)

	s.mailer.On("ConfirmationMail", mock.Anything, &templates.ConfirmationReq{
		Mail:          email,
		ParentName:    pname,
		ParentSurname: psurname,
		EventName:     event.Title,
		Name:          name,
		Surname:       surname,
		Pills:         "-",
		Restrictions:  "-",
		Info:          "",
		PhotoURL:      event.OwnerPhoto,
		Sum:           day.Price - event.PromoDiscount,
		Owner:         "John Doe",
		Text:          event.OwnerPhone + " " + event.OwnerEmail,
		Days:          []string{day.Description},
		RegInfo:       *event.Info,
		Payment: templates.PaymentDetails{
			IBAN:             event.IBAN,
			PaymentReference: event.PaymentReference,
			SpecificSymbol:   s.expectedNextSpecificSymbol(),
			Link:             s.expectedPayMeLink(day.Price-event.PromoDiscount, event.IBAN, name, surname, event, s.expectedNextSpecificSymbol()),
			QRCode:           "",
		},
	}).Return(nil)

	createRegistrationReq := func() (*http.Request, *httptest.ResponseRecorder) {
		u := fmt.Sprintf("/api/registrations/%d", event.ID)
		return s.NewRequest(http.MethodPost, u, echo.Map{
			"child": echo.Map{
				"name":               name,
				"surname":            surname,
				"gender":             gender,
				"city":               city,
				"finishedSchoolYear": school,
				"dateOfBirth":        birth,
			},
			"parent": echo.Map{
				"name":    pname,
				"surname": psurname,
				"email":   email,
				"phone":   phone,
			},
			"days":       []interface{}{day.ID},
			"promo_code": promoCode,
		})
	}

	// Promo registration not active
	_, err := s.dbx.Exec("UPDATE events SET active = false, promo_registration = false WHERE id = $1", event.ID)
	s.Require().NoError(err)

	req, rec = createRegistrationReq()
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{"event_id": "not active"},
		}, body)
	})

	// Promo registration active
	_, err = s.dbx.Exec("UPDATE events SET active = false, promo_registration = true WHERE id = $1", event.ID)
	s.Require().NoError(err)

	req, rec = createRegistrationReq()
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.NotEmpty(body["token"])
		delete(body, "token")
		s.Equal(echo.Map{
			"registeredIDs": []interface{}{float64(day.ID)},
			"success":       true,
		}, body)
	})
	// TODO: assert promo code

	// Send another request with the same promo code
	req, rec = createRegistrationReq()
	s.AssertServerResponseObject(req, rec, http.StatusBadRequest, func(body echo.Map) {
		s.Equal(echo.Map{
			"error": "Token already used.",
		}, body)
	})
}

func (s *RegistrationSuite) TestRegister_NotActive_Promo_JWT() {
	s.testRegister_NotActive_Promo()
}

func (s *RegistrationSuite) TestRegister_NotActive_Promo_Simple() {
	s.promoRegistry = promo.NewRegistry(s.postgresDB, promo.NewSimpleGenerator(s.logger), s.mailer)
	s.testRegister_NotActive_Promo()
}

func (s *RegistrationSuite) TestRegister_OK() {
	const (
		name     = "dano"
		surname  = "Žščťžzharmanca"
		pname    = "janko"
		psurname = "hrasko"
		gender   = "male"
		city     = "BB"
		phone    = "+421"
		email    = "dano@mail.sk"
		school   = "3.ZS"
	)
	event := s.InsertEvent()

	day := event.Days[0]

	birth := time.Now().Format(time.RFC3339)

	u := fmt.Sprintf("/api/registrations/%d", event.ID)
	req, rec := s.NewRequest(http.MethodPost, u, echo.Map{
		"child": echo.Map{
			"name":               name,
			"surname":            surname,
			"gender":             gender,
			"city":               city,
			"finishedSchoolYear": school,
			"dateOfBirth":        birth,
		},
		"parent": echo.Map{
			"name":    pname,
			"surname": psurname,
			"email":   email,
			"phone":   phone,
		},
		"days": []interface{}{day.ID},
	})
	s.mailer.On("ConfirmationMail", mock.Anything, &templates.ConfirmationReq{
		Mail:          email,
		ParentName:    pname,
		ParentSurname: psurname,
		EventName:     event.Title,
		Name:          name,
		Surname:       surname,
		Pills:         "-",
		Restrictions:  "-",
		Info:          "",
		PhotoURL:      event.OwnerPhoto,
		Sum:           day.Price,
		Owner:         "John Doe",
		Text:          event.OwnerPhone + " " + event.OwnerEmail,
		Days:          []string{day.Description},
		RegInfo:       *event.Info,
		Payment: templates.PaymentDetails{
			IBAN:             event.IBAN,
			PaymentReference: event.PaymentReference,
			SpecificSymbol:   s.expectedNextSpecificSymbol(),
			Link:             s.expectedPayMeLink(day.Price, event.IBAN, name, surname, event, s.expectedNextSpecificSymbol()),
			QRCode:           "",
		},
	}).Return(nil)

	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.NotEmpty(body["token"])
		delete(body, "token")
		s.Equal(echo.Map{
			"registeredIDs": []interface{}{float64(day.ID)},
			"success":       true,
		}, body)
	})
	//TODO: assert DB content
	s.mailer.AssertExpectations(s.T())
}

func (s *RegistrationSuite) expectedNextSpecificSymbol() string {
	specificSymbol := s.lastSequenceValue("specific_symbol_seq")
	return strconv.Itoa(specificSymbol + 1)
}

func (s *RegistrationSuite) TestDelete_Unauthorized() {
	req, rec := s.NewRequest(http.MethodDelete, "/api/registrations/42", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *RegistrationSuite) TestDelete_NotFound() {
	req, rec := s.NewRequest(http.MethodDelete, "/api/registrations/42", nil)
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, nil)
}

func (s *RegistrationSuite) TestDelete_OK() {
	event := s.InsertEvent()
	reg := s.createRegistration(event.Days[0].ID)

	u := fmt.Sprintf("/api/registrations/%d", reg.ID)
	req, rec := s.NewRequest(http.MethodDelete, u, nil)
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusOK, nil)
}

func (s *RegistrationSuite) expectedPayMeLink(
	amount int,
	iban string,
	name string,
	surname string,
	event *model.Event,
	specificSymbol string,
) string {
	v := url.Values{}
	v.Set("AM", strconv.Itoa(amount))
	v.Set("CC", "EUR")
	v.Set("V", "1")
	v.Set("CN", "salezko")
	v.Set("IBAN", iban)
	v.Set("PI", fmt.Sprintf("/VS%s/SS%s/KS%s", event.PaymentReference, specificSymbol, ""))
	v.Set("MSG", fmt.Sprintf("%s %s %s", event.Title, name, surname))
	return "https://payme.sk?" + v.Encode()
}

func (s *RegistrationSuite) TestResendConfirmationEmail_Unauthorized() {
	req, rec := s.NewRequest(http.MethodPost, "/api/registrations/42/resend_confirmation", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *RegistrationSuite) TestResendConfirmationEmail_InvalidEmail() {
	req, rec := s.NewRequest(http.MethodPost, "/api/registrations/42/resend_confirmation", nil)
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{"email": "invalid"},
		}, body)
	})
}

func (s *RegistrationSuite) TestResendConfirmationEmail_NotFound() {
	req, rec := s.NewRequest(http.MethodPost, "/api/registrations/42/resend_confirmation", echo.Map{"email": "a@b.com"})
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, nil)
}

func (s *RegistrationSuite) TestResendConfirmationEmail_Success() {
	event := s.InsertEvent()
	reg := s.createRegistration(event.Days[0].ID)

	const newEmail = "a@b.com"

	s.mailer.On("ConfirmationMail", mock.Anything, &templates.ConfirmationReq{
		Mail:          newEmail,
		ParentName:    reg.ParentName,
		ParentSurname: reg.ParentSurname,
		EventName:     event.Title,
		Name:          reg.Name,
		Surname:       reg.Surname,
		Pills:         "-",
		Restrictions:  "-",
		Info:          "",
		PhotoURL:      event.OwnerPhoto,
		Sum:           reg.AmountToPay(),
		Owner:         "John Doe",
		Text:          event.OwnerPhone + " " + event.OwnerEmail,
		Days:          []string{event.Days[0].Description},
		RegInfo:       *event.Info,
		Payment: templates.PaymentDetails{
			IBAN:             event.IBAN,
			PaymentReference: event.PaymentReference,
			SpecificSymbol:   reg.SpecificSymbol,
			Link:             s.expectedPayMeLink(reg.AmountToPay(), event.IBAN, reg.Name, reg.Surname, event, reg.SpecificSymbol),
			QRCode:           "",
		},
	}).Return(nil)

	u := fmt.Sprintf("/api/registrations/%d/resend_confirmation", reg.ID)
	req, rec := s.NewRequest(http.MethodPost, u, echo.Map{"email": newEmail})
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusAccepted, nil)
	s.mailer.AssertExpectations(s.T())

	updatedRegistration, err := repository.FindRegistrationByID(context.Background(), s.dbx, reg.ID)
	s.Require().NoError(err)
	s.Require().Equal(newEmail, updatedRegistration.Email)
}

func (s *RegistrationSuite) TestSendPaymentReminder_Unauthorized() {
	req, rec := s.NewRequest(http.MethodPost, "/api/send_payment_reminder", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *RegistrationSuite) TestSendPaymentReminder_EventNotFound() {
	req, rec := s.NewRequest(http.MethodPost, "/api/send_payment_reminder",
		echo.Map{"event_id": 42})
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, nil)
}

func (s *RegistrationSuite) TestSendPaymentReminder_SkipPayed() {
	event := s.InsertEvent()
	_ = s.createRegistration(event.Days[0].ID, func(registration *model.Registration) {
		registration.Payed = s.intRef(registration.AmountToPay())
	})

	req, rec := s.NewRequest(http.MethodPost, "/api/send_payment_reminder", echo.Map{"event_id": event.ID})
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"sent": float64(0), "finished_all": true}, body)
	})
	s.mailer.AssertExpectations(s.T())
}

func (s *RegistrationSuite) TestSendPaymentReminder_Success() {
	event := s.InsertEvent()
	reg := s.createRegistration(event.Days[0].ID)

	s.mailer.On("PaymentReminderMail", mock.Anything, &templates.PaymentReminderReq{
		Mail:          reg.Email,
		ParentName:    reg.ParentName,
		ParentSurname: reg.ParentSurname,
		EventName:     event.Title,
		Name:          reg.Name,
		Surname:       reg.Surname,
		Sum:           reg.AmountToPay(),
		Payment: templates.PaymentDetails{
			IBAN:             event.IBAN,
			PaymentReference: event.PaymentReference,
			SpecificSymbol:   reg.SpecificSymbol,
			Link:             s.expectedPayMeLink(reg.AmountToPay(), event.IBAN, reg.Name, reg.Surname, event, reg.SpecificSymbol),
			QRCode:           "",
		},
	}).Return(nil)

	req, rec := s.NewRequest(http.MethodPost, "/api/send_payment_reminder", echo.Map{"event_id": event.ID})
	s.AuthorizeRequest(req, &auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "admin@sbb.sk",
		},
	})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"sent": float64(1), "finished_all": true}, body)
	})
	s.mailer.AssertExpectations(s.T())
}

func TestRegistrationSuite(t *testing.T) {
	suite.Run(t, new(RegistrationSuite))
}
