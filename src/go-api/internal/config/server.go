package config

type Server struct {
	Host           string   `envconfig:"HOST" default:"0.0.0.0"`
	Port           int      `envconfig:"PORT" default:"5000"`
	JWTSecret      []byte   `envconfig:"JWT_SECRET" required:"true"`
	AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"https://leto-admin.salezko.sk,https://leto.salezko.sk"`
	CookieDomain   string   `envconfig:"COOKIE_DOMAIN" default:"salezko.sk"`
	CookieSecure   bool     `envconfig:"COOKIE_SECURE" default:"true"`
}
