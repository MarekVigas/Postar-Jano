package api_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/random"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/MarekVigas/Postar-Jano/internal/auth"
	"github.com/MarekVigas/Postar-Jano/internal/config"
	"github.com/MarekVigas/Postar-Jano/internal/promo"

	"github.com/MarekVigas/Postar-Jano/internal/api"
	"github.com/MarekVigas/Postar-Jano/internal/mailer/templates"
	"github.com/MarekVigas/Postar-Jano/internal/model"
	"github.com/MarekVigas/Postar-Jano/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	loggingEnabled = true
	jwtSecret      = "top-secret"
	promoSecret    = "secret"
)

type CommonSuite struct {
	suite.Suite
	logger *zap.Logger

	api          *api.API
	rootDB       *sql.DB
	db           *sql.DB
	dbx          *sqlx.DB
	mailer       *SenderMock
	promoManager repository.PromoManager
	dbName       string
}

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) ConfirmationMail(ctx context.Context, req *templates.ConfirmationReq) error {
	return m.Called(ctx, req).Error(0)
}

func (m *SenderMock) PromoMail(ctx context.Context, req *templates.PromoReq) error {
	return m.Called(ctx, req).Error(0)
}

func (m *SenderMock) NotificationMail(ctx context.Context, req *templates.NotificationReq) error {
	return m.Called(ctx, req).Error(0)
}

func (s *CommonSuite) SetupSuite() {
	var err error

	if loggingEnabled {
		s.logger, err = zap.NewDevelopment()
		s.Require().NoError(err)
	} else {
		s.logger = zap.NewNop()
	}

	var dbConfig config.DB
	s.Require().NoError(envconfig.Process("", &dbConfig))

	s.rootDB, err = dbConfig.Connect()
	s.Require().NoError(err)

	s.dbName = fmt.Sprintf("testing" + random.String(10, random.Lowercase))

	// Create db schema.
	_, err = s.rootDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, s.dbName))
	s.Require().NoError(err)
	_, err = s.rootDB.Exec(fmt.Sprintf(`CREATE DATABASE %s;`, s.dbName))
	s.Require().NoError(err)

	dbConfig.Database = s.dbName

	err = repository.RunMigrations(s.logger, &dbConfig, s.dbName, "file://../../migrations")
	s.Require().NoError(err)

	s.db, err = dbConfig.Connect()
	s.Require().NoError(err)
	s.dbx = sqlx.NewDb(s.db, "postgres")

	s.promoManager = promo.NewJWTGenerator(s.logger, []byte(promoSecret), nil, nil)

}

func (s *CommonSuite) TearDownSuite() {
	_ = s.db.Close()
	_, _ = s.rootDB.Exec("DROP DATABASE " + s.dbName)
	_ = s.rootDB.Close()
	_ = s.logger.Sync()
}

func (s *CommonSuite) SetupTest() {
	ctx := context.Background()

	s.mailer = &SenderMock{}

	repo := repository.NewPostgresRepo(s.db, s.promoManager)
	s.api = api.New(
		s.logger,
		repo,
		auth.NewFromDB(repo),
		s.mailer,
		[]byte(jwtSecret),
	)
	s.NoError(repository.Reset(ctx, s.db))
}

func (s *CommonSuite) AssertServerResponseObject(
	req *http.Request,
	rec *httptest.ResponseRecorder,
	expectedStatus int,
	checkBody func(echo.Map),
) {
	s.api.ServeHTTP(rec, req)
	if !s.Equal(expectedStatus, rec.Code) {
		s.T().Log("Response body:\n", rec.Body.String())
		return
	}
	if checkBody != nil {
		var body echo.Map
		if s.NoError(json.NewDecoder(rec.Body).Decode(&body)) {
			checkBody(body)
		}
	}
}

func (s *CommonSuite) NewRequest(
	method string,
	target string,
	body interface{},
) (*http.Request, *httptest.ResponseRecorder) {

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if s.NoError(err) {
			bodyReader = bytes.NewReader(b)
		}
	}

	req := httptest.NewRequest(method, target, bodyReader)
	if body != nil {
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return req, rec
}

func (s *CommonSuite) AssertServerResponseArray(
	req *http.Request,
	rec *httptest.ResponseRecorder,
	expectedStatus int,
	checkBody func([]interface{}),
) {
	s.api.ServeHTTP(rec, req)
	if !s.Equal(expectedStatus, rec.Code) {
		s.T().Log("Response body:\n", rec.Body.String())
		return
	}
	if checkBody != nil {
		var body []interface{}
		if s.NoError(json.NewDecoder(rec.Body).Decode(&body)) {
			checkBody(body)
		}
	}
}

func (s *CommonSuite) InsertEvent() *model.Event {
	ctx := context.Background()
	owner, err := (&model.Owner{
		Name:     "John",
		Surname:  "Doe",
		Email:    "john@doe.com",
		Username: "john@doe.com",
		Phone:    "123",
		Photo:    "phot.jpg",
		Gender:   "M",
	}).Create(ctx, s.dbx, "bla bla")
	s.Require().NoError(err)

	event := model.Event{
		Title:            "Camp 42",
		Description:      "Lorem ipsum",
		DateFrom:         "15 june 2020",
		DateTo:           "20 june 2020",
		Location:         "somewhere",
		MinAge:           10,
		MaxAge:           15,
		Info:             s.stringRef("xyz.."),
		Photo:            "photo",
		Active:           true,
		IBAN:             "SK7409000000005073149693",
		PromoDiscount:    2,
		PaymentReference: "202401",
		EventOwner:       s.eventOwnerFromModel(owner),
	}
	s.Require().NoError((&event).Create(ctx, s.dbx))

	day := model.Day{
		Description: "Desc",
		Capacity:    10,
		LimitBoys:   s.intRef(5),
		LimitGirls:  s.intRef(5),
		Price:       12,
		EventID:     event.ID,
	}
	s.Require().NoError((&day).Create(ctx, s.dbx))

	event.Days = append(event.Days, day)

	return &event
}

func (s *CommonSuite) createRegistration() *model.Registration {
	reg, err := (&model.Registration{
		Name:               "sdafa",
		Surname:            "asdfasf",
		Gender:             "female",
		DateOfBirth:        time.Time{},
		FinishedSchool:     "1ZS",
		AttendedPrevious:   true,
		AttendedActivities: nil,
		City:               "fadsf",
		Pills:              nil,
		Problems:           nil,
		Notes:              "",
		ParentName:         "sadfa",
		ParentSurname:      "asdfa",
		Email:              "",
		Phone:              "",
		Amount:             0,
		Payed:              nil,
		Discount:           nil,
		AdminNote:          "asdfas",
		Token:              "asfasfd",
	}).Create(context.Background(), s.dbx)
	s.Require().NoError(err)

	return reg
}

func (s *CommonSuite) stringRef(str string) *string {
	return &str
}

func (s *CommonSuite) intRef(val int) *int {
	return &val
}

func (s *CommonSuite) lastSequenceValue(seqName string) int {
	row := s.db.QueryRow(fmt.Sprintf("SELECT last_value FROM %s", seqName))
	s.Require().NoError(row.Err())

	var value int
	s.Require().NoError(row.Scan(&value))
	return value
}

func (s *CommonSuite) eventOwnerFromModel(owner *model.Owner) model.EventOwner {
	return model.EventOwner{
		OwnerID:      owner.ID,
		OwnerName:    owner.Name,
		OwnerSurname: owner.Surname,
		OwnerEmail:   owner.Email,
		OwnerPhone:   owner.Phone,
		OwnerPhoto:   owner.Photo,
		OwnerGender:  owner.Gender,
	}
}
