package api_test

import (
	"net/http"
	"testing"

	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type AuthSuite struct {
	CommonSuite
}

func (s *AuthSuite) TestSignIn_UnprocessableEntity_EmptyUsername() {
	req, rec := s.NewRequest(http.MethodPost, "/api/sign/in", echo.Map{
		"password": "xyz",
	})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{"errors": map[string]interface{}{
			"username": "required",
		}}, body)
	})
}

func (s *AuthSuite) TestSignIn_UnprocessableEntity_EmptyPass() {
	req, rec := s.NewRequest(http.MethodPost, "/api/sign/in", echo.Map{
		"username": "xyz",
	})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{"errors": map[string]interface{}{
			"password": "required",
		}}, body)
	})
}

func (s *AuthSuite) TestSign_Forbidden() {
	const (
		userName = "test"
		pass     = "nbusr123"
	)
	_ = s.addOwner(userName, pass)
	req, rec := s.NewRequest(http.MethodPost, "/api/sign/in", echo.Map{
		"username": userName,
		"password": "INCORRECT_PASS",
	})
	s.AssertServerResponseObject(req, rec, http.StatusForbidden, nil)
}

func (s *AuthSuite) TestSignIn_OK() {
	const (
		userName = "test"
		pass     = "nbusr123"
	)
	user := s.addOwner(userName, pass)
	req, rec := s.NewRequest(http.MethodPost, "/api/sign/in", echo.Map{
		"username": userName,
		"password": pass,
	})
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.NotEmpty(body["token"])
		var claims auth.Claims
		token, err := jwt.ParseWithClaims(body["token"].(string), &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		s.NoError(err)
		s.NotNil(token)
		s.Equal(claims.Id, user.Email)
	})
}

func (s *AuthSuite) addOwner(username string, pass string) *model.Owner {
	psswd, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	s.Require().NoError(err)

	owner := model.Owner{
		ID:       0,
		Name:     random.String(10, random.Lowercase),
		Surname:  random.String(10, random.Lowercase),
		Email:    username + "example.com",
		Username: username,
		Pass:     string(psswd),
		Phone:    "0908",
		Photo:    "xyz",
		Gender:   "male",
	}
	_, err = sqlx.NamedExec(s.dbx, `INSERT INTO owners(
			name,
			surname,
			email,
			username,
			pass,
			phone,
			photo,
			gender
		) VALUES (
			:name,
			:surname,
			:email,
			:username,
			:pass,
			:phone,
			:photo,
			:gender
		)`, &owner)
	s.NoError(err)
	return &owner
}

func (s *AuthSuite) TestListRegistrations_Unauthorized() {
	req, rec := s.NewRequest(http.MethodGet, "/api/registrations", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *AuthSuite) TestListRegistrations_OK_Empty() {
	req, rec := s.NewRequest(http.MethodGet, "/api/registrations", nil)
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusOK, nil)
}

func (s *AuthSuite) TestPutRegistrations_Unauthorized() {
	req, rec := s.NewRequest(http.MethodPut, "/api/registrations/1", nil)
	s.AssertServerResponseObject(req, rec, http.StatusUnauthorized, nil)
}

func (s *AuthSuite) TestPutRegistrations_NotFound() {
	req, rec := s.NewRequest(http.MethodPut, "/api/registrations/42", echo.Map{
		"child": echo.Map{
			"name":    "john",
			"surname": "doe",
		},
		"parent": echo.Map{"email": "me@example.com"},
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, nil)
}

func (s *AuthSuite) TestPutRegistrations_UnprocessableEntity_MissingChildName() {
	s.T().Skipf("Child name is not updated")
	req, rec := s.NewRequest(http.MethodPut, "/api/registrations/42", echo.Map{
		"child": echo.Map{
			"surname": "doe",
		},
		"parent": echo.Map{"email": "me@example.com"},
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{
				"updatereq.child.name": "required",
			},
		}, body)
	})
}

func (s *AuthSuite) TestPutRegistrations_UnprocessableEntity_InvalidMail() {
	s.T().Skipf("Email is not updated")
	req, rec := s.NewRequest(http.MethodPut, "/api/registrations/42", echo.Map{
		"child": echo.Map{
			"name":    "john",
			"surname": "doe",
		},
		"parent": echo.Map{"email": "bla"},
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusUnprocessableEntity, func(body echo.Map) {
		s.Equal(echo.Map{
			"errors": map[string]interface{}{
				"updatereq.parent.email": "email",
			},
		}, body)
	})
}

func (s *AuthSuite) TestPutRegistrations_OK() {
	event := s.InsertEvent()

	_, err := s.db.Exec(`INSERT INTO days (
		id,
		capacity,
		limit_boys,
		limit_girls,
		description,
		price,
		event_id
	) VALUES (
		12,
		10,
		5,
		5,
		'bla',
		42,
		$1
	)`, event.ID)
	s.Require().NoError(err)

	_, err = s.db.Exec(`INSERT INTO registrations(
		id,
		name,
		surname,
		token,
		gender,
		amount,
		payed,
		finished_school,
		attended_previous,
		city,
		pills,
		notes,
		parent_name,
		parent_surname,
		email,
		phone,
		date_of_birth,
		created_at,
		updated_at
	) VALUES (
		15,
		'sadf',
		'sadf',
		'sadf',
		'female',
		10,
		0,
		'zs',
		true,
		'bb',
		'pills',
		'notest',
		'parentN',
		'parentS',
		'email',
		'phone',
		NOW(),
		NOW(),
		NOW()
	)`)
	s.Require().NoError(err)

	_, err = s.db.Exec(`INSERT INTO signups(
		day_id,
		registration_id,
		state,
		created_at,
		updated_at
	) VALUES (
		12,
		15,
		'sadf',
		NOW(),
		NOW()
	)`)
	s.Require().NoError(err)

	req, rec := s.NewRequest(http.MethodPut, "/api/registrations/15", echo.Map{
		"child": echo.Map{
			"name":    "john",
			"surname": "doe",
		},
		"parent":     echo.Map{"email": "me@example.com"},
		"payed":      nil,
		"discount":   nil,
		"admin_note": nil,
	})
	s.AuthorizeRequest(req, &auth.Claims{})
	s.AssertServerResponseObject(req, rec, http.StatusAccepted, nil)
	//TODO: test db content
}

func (s *CommonSuite) AuthorizeRequest(req *http.Request, claims *auth.Claims) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := tok.SignedString([]byte(jwtSecret))
	s.NoError(err)

	req.Header.Set(echo.HeaderAuthorization, "Bearer "+ss)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
