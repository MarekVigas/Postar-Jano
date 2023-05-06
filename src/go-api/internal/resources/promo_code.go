package resources

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PromoCodeReq struct {
	Email             string `json:"email"   validate:"email,required"`
	RegistrationCount int    `json:"registration_count" validate:"required"`
	SendEmail         bool   `json:"send_email"`
}

func (r *PromoCodeReq) Validate() interface{} {
	v := validator.New()
	if errs := v.Struct(r); errs != nil {
		msgs := echo.Map{}
		if e, ok := errs.(validator.ValidationErrors); ok {
			for _, err := range e {
				msgs[strings.ToLower(err.StructNamespace())] = err.Tag()
			}
		}
		return echo.Map{"errors": msgs}
	}
	return nil
}

type ValidatePromoCodeReq struct {
	PromoCode string `json:"promo_code"   validate:"required"`
}

func (r *ValidatePromoCodeReq) Validate() interface{} {
	v := validator.New()
	if errs := v.Struct(r); errs != nil {
		msgs := echo.Map{}
		if e, ok := errs.(validator.ValidationErrors); ok {
			for _, err := range e {
				msgs[strings.ToLower(err.StructNamespace())] = err.Tag()
			}
		}
		return echo.Map{"errors": msgs}
	}
	return nil
}
