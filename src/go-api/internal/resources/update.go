package resources

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UpdateReq struct {
	Child struct {
		Name             string    `json:"name" validate:"required"`
		Surname          string    `json:"surname" validate:"required"`
		Gender           string    `json:"gender"`
		DateOfBirth      time.Time `json:"dateOfBirth"`
		FinishedSchool   string    `json:"finishedSchoolYear"`
		AttendedPrevious bool      `json:"attendedPreviousEvents"`
		City             string    `json:"city"`
	} `json:"child"`

	Medicine struct {
		Pills string `json:"drugs"`
	} `json:"medicine"`

	Health struct {
		Problems string `json:"problems"`
	} `json:"health"`

	Parent struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Email   string `json:"email" validate:"email,required"`
		Phone   string `json:"phone"`
	} `json:"parent"`

	Membership struct {
		AttendedActivities string `json:"attendedActivities"`
	} `json:"memberShip"`
	Notes     string `json:"notes"`
	Amount    int    `json:"amount"`
	Payed     *int   `json:"payed"`
	Discount  *int   `json:"discount"`
	AdminNote string `json:"admin_note"`
}

func (s *UpdateReq) Validate() interface{} {
	v := validator.New()
	if errs := v.Struct(s); errs != nil {
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
