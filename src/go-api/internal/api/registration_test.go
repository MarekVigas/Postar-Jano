package api_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"

	"github.com/dgrijalva/jwt-go"
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

func (s *RegistrationSuite) TestRegister_NotActive_Promo() {
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

	_, err := s.dbx.Exec("UPDATE events SET active = false, promo_registration = true WHERE id = $1", event.ID)
	s.Require().NoError(err)

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
	u := fmt.Sprintf("/api/registrations/%d", event.ID)
	req, rec = s.NewRequest(http.MethodPost, u, echo.Map{
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
	}).Return(nil)

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
	req, rec = s.NewRequest(http.MethodPost, u, echo.Map{
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

	s.AssertServerResponseObject(req, rec, http.StatusBadRequest, func(body echo.Map) {
		s.Equal(echo.Map{
			"error": "Token already used.",
		}, body)
	})
}

func (s *RegistrationSuite) TestRegister_OK() {
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
}

func (s *RegistrationSuite) TestDelete_Unauthorized() {
	req, rec := s.NewRequest(http.MethodDelete, "/api/registrations/42", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *RegistrationSuite) TestDelete_NotFound() {
	req, rec := s.NewRequest(http.MethodDelete, "/api/registrations/42", nil)
	s.AuthorizeRequest(req, &auth.Claims{StandardClaims: jwt.StandardClaims{Id: "admin@sbb.sk"}})
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, nil)
}

func (s *RegistrationSuite) TestDelete_OK() {
	reg := s.createRegistration()

	u := fmt.Sprintf("/api/registrations/%d", reg.ID)
	req, rec := s.NewRequest(http.MethodDelete, u, nil)
	s.AuthorizeRequest(req, &auth.Claims{StandardClaims: jwt.StandardClaims{Id: "admin@sbb.sk"}})
	s.AssertServerResponseObject(req, rec, http.StatusOK, nil)
}

func TestRegistrationSuite(t *testing.T) {
	suite.Run(t, new(RegistrationSuite))
}
