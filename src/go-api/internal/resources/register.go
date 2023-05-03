package resources

import "time"

type RegisterReq struct {
	Child struct {
		Name             string    `json:"name"    validate:"required"`
		Surname          string    `json:"surname" validate:"required"`
		Gender           string    `json:"gender"  validate:"eq=male|eq=female"`
		DateOfBirth      time.Time `json:"dateOfBirth" validate:"required"`
		FinishedSchool   string    `json:"finishedSchoolYear" validate:"required"`
		AttendedPrevious bool      `json:"attendedPreviousEvents"`
		City             string    `json:"city"      validate:"required"`
	} `json:"child"`

	Medicine struct {
		Pills *string `json:"drugs"`
	} `json:"medicine"`

	Health struct {
		Problems *string `json:"problems"`
	} `json:"health"`

	Parent struct {
		Name    string `json:"name"    validate:"required"`
		Surname string `json:"surname" validate:"required"`
		Email   string `json:"email"   validate:"email,required"`
		Phone   string `json:"phone"   validate:"required"`
	} `json:"parent"`

	Membership struct {
		AttendedActivities *string `json:"attendedActivities"`
	} `json:"memberShip"`
	Notes     string  `json:"notes"`
	DayIDs    []int   `json:"days" validate:"required"`
	PromoCode *string `json:"promo_code"`
}
