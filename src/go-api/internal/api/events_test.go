package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/MarekVigas/Postar-Jano/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type EventsSuite struct {
	CommonSuite
}

func (s *EventsSuite) TestListEvents_OK() {
	event := s.InsertEvent()

	req, rec := s.NewRequest(http.MethodGet, "/api/events", nil)
	s.AssertServerResponseArray(req, rec, http.StatusOK, func(body []interface{}) {
		s.Equal([]interface{}{
			map[string]interface{}(s.eventToResource(event)),
		}, body)
	})
}

func (s *EventsSuite) TestFindEvent_NotFound() {
	req, rec := s.NewRequest(http.MethodGet, "/api/events/1", nil)
	s.AssertServerResponseObject(req, rec, http.StatusNotFound, func(body echo.Map) {
		s.Equal(echo.Map{"message": "Not Found"}, body)
	})
}

func (s *EventsSuite) TestFindEvent_OK() {
	event := s.InsertEvent()

	u := fmt.Sprintf("/api/events/%d", event.ID)
	req, rec := s.NewRequest(http.MethodGet, u, nil)
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(s.eventToResource(event), body)
	})
}

func (s *EventsSuite) eventToResource(event *model.Event) echo.Map {
	res := echo.Map{
		"id":          float64(event.ID),
		"title":       event.Title,
		"description": event.Description,
		"date_from":   event.DateFrom,
		"date_to":     event.DateTo,
		"location":    event.Location,
		"min_age":     float64(event.MinAge),
		"max_age":     float64(event.MaxAge),
		"owner": map[string]interface{}{
			"id":      float64(event.OwnerID),
			"email":   event.OwnerEmail,
			"name":    event.OwnerName,
			"surname": event.OwnerSurname,
			"phone":   event.OwnerPhone,
			"photo":   event.OwnerPhoto,
			"gender":  event.OwnerGender,
		},

		"active":             event.Active,
		"promo_registration": event.PromoRegistration,

		"photo":     event.Photo,
		"info":      nil,
		"mail_info": nil,
		"price":     nil,
		"time":      nil,
	}
	if event.Time != nil {
		res["time"] = *event.Time
	}
	if event.Price != nil {
		res["price"] = *event.Price
	}
	if event.Info != nil {
		res["info"] = *event.Info
	}
	if event.MailInfo != nil {
		res["mail_info"] = *event.MailInfo
	}
	days := make([]interface{}, len(event.Days))
	for i, d := range event.Days {
		day := map[string]interface{}{
			"id":          float64(d.ID),
			"capacity":    float64(d.Capacity),
			"description": d.Description,
			"price":       float64(d.Price),
			"limit_boys":  nil,
			"limit_girls": nil,
		}
		if d.LimitBoys != nil {
			day["limit_boys"] = float64(*d.LimitBoys)
		}
		if d.LimitGirls != nil {
			day["limit_girls"] = float64(*d.LimitGirls)
		}
		days[i] = day
	}
	res["days"] = days
	return res
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}
