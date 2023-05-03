package config

import "time"

type Promo struct {
	Secret         []byte     `envconfig:"PROMO_SECRET" required:"true"`
	ExpirationDate *time.Time `envconfig:"PROMO_EXPIRATION_DATE"`
	ActivationDate *time.Time `envconfig:"PROMO_ACTIVATION_DATE"`
}
