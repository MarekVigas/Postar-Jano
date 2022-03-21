package resources

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *SignIn) Validate() interface{} {
	v := validator.New()
	if errs := v.Struct(s); errs != nil {
		msgs := echo.Map{}
		if e, ok := errs.(validator.ValidationErrors); ok {
			for _, err := range e {
				msgs[strings.ToLower(err.Field())] = err.Tag()
			}
		}
		return echo.Map{"errors": msgs}
	}
	return nil
}
