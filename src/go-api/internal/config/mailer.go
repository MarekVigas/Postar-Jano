package config

type Mailer struct {
	MailgunDomain            string `envconfig:"MAILGUN_DOMAIN" required:"true"`
	MailgunKey               string `envconfig:"MAILGUN_KEY"    required:"true"`
	EUBase                   bool   `envconfig:"MAILGUN_EU_BASE" required:"true"`
	ConfirmationMailTemplate string `envconfig:"CONFIRMATION_EMAIL_TEMPLATE" required:"true"`
	PromoMailTemplate        string `envconfig:"PROMO_EMAIL_TEMPLATE" required:"true"`
	NotificationMailTemplate string `envconfig:"NOTIFICATION_EMAIL_TEMPLATE" required:"false"`
}
