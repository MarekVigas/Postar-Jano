package resources

type PromoCodeReq struct {
	Email             string `json:"email"   validate:"email,required"`
	RegistrationCount int    `json:"registration_count" validate:"required"`
	SendEmail         bool   `json:"send_email"`
}

type ValidatePromoCodeReq struct {
	PromoCode string `json:"promo_code"   validate:"required"`
}

type ValidatePromoCodeResp struct {
	Status                 string `json:"status"`
	AvailableRegistrations int    `json:"available_registrations"`
}

func ValidPromoCodeResp(availableRegistrations int) *ValidatePromoCodeResp {
	return &ValidatePromoCodeResp{
		Status:                 "ok",
		AvailableRegistrations: availableRegistrations,
	}
}

func InvalidPromoCodeResp() *ValidatePromoCodeResp {
	return &ValidatePromoCodeResp{Status: "invalid"}
}
