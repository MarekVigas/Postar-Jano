package model

type Event struct {
	ID                int     `json:"id" db:"id"`
	Title             string  `json:"title" db:"title"`
	Description       string  `json:"description" db:"description"`
	DateFrom          string  `json:"date_from" db:"date_from"`
	DateTo            string  `json:"date_to" db:"date_to"`
	Location          string  `json:"location" db:"location"`
	MinAge            int     `json:"min_age" db:"min_age"`
	MaxAge            int     `json:"max_age" db:"max_age"`
	Info              *string `json:"info" db:"info"`
	Photo             string  `json:"photo" db:"photo"`
	Time              *string `json:"time" db:"time"`
	Price             *string `json:"price" db:"price"`
	MailInfo          *string `json:"mail_info" db:"mail_info"`
	Active            bool    `json:"active" db:"active"`
	PromoRegistration bool    `json:"promo_registration" db:"promo_registration"`
	IBAN              string  `json:"iban" db:"iban"`
	PromoDiscount     int     `json:"promo_discount" db:"promo_discount"`
	PaymentReference  string  `json:"payment_reference" db:"payment_reference"`
	// Refactor embeded struct
	EventOwner `json:"owner"`
	Days       []Day `json:"days"`
}

type EventOwner struct {
	OwnerID      int    `json:"id" db:"owner_id"`
	OwnerName    string `json:"name" db:"owner_name"`
	OwnerSurname string `json:"surname" db:"owner_surname"`
	OwnerEmail   string `json:"email" db:"owner_email"`
	OwnerPhone   string `json:"phone" db:"owner_phone"`
	OwnerPhoto   string `json:"photo" db:"owner_photo"`
	OwnerGender  string `json:"gender" db:"owner_gender"`
}
