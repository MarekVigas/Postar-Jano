package config

type Server struct {
	Host      string `envconfig:"HOST" default:"0.0.0.0"`
	Port      int    `envconfig:"PORT" default:"5000"`
	JWTSecret []byte `envconfig:"JWT_SECRET" required:"true"`
}
