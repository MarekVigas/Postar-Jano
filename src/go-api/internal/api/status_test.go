package api_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/suite"
)

type StatusSuite struct {
	CommonSuite
}

func (s *StatusSuite) TestStatus_OK() {
	req, rec := s.NewRequest(http.MethodGet, "/api/status", nil)
	s.AssertServerResponseObject(req, rec, http.StatusOK, func(body echo.Map) {
		s.Equal(echo.Map{"status": "ok"}, body)
	})
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}
